package digest

import (
	"fmt"
	"strings"
	"sync"

	"github.com/blainemoser/errorLogging/files"
	"github.com/blainemoser/errorLogging/filewatcher"
	"github.com/blainemoser/errorLogging/slack"
	"github.com/blainemoser/errorLogging/utils"
	"github.com/fsnotify/fsnotify"
)

type queue struct {
	events     map[int]*filewatcher.Event
	lock       sync.RWMutex
	processing bool
	slack      *slack.Slack
}

var q *queue

func Start() {
	q = &queue{
		events: make(map[int]*filewatcher.Event),
		slack:  slack.NewSlack(),
	}
}

func AddToQueue(event fsnotify.Event) error {
	if event.Name == "" {
		return nil
	}
	err := files.List.Reset(event.Name)
	if err != nil {
		return err
	}
	var index int
	if checkWrite(event) {
		file, err := files.List.GetFile(event.Name)
		if err != nil {
			return err
		}
		file.Watcher.Start()
		e := file.Watcher.GetEvent()
		if file.Watcher.AppendedTo() && !file.Filtered(e) && !file.Suppressed(e) {
			q.lock.Lock()
			index = len(q.events) + 1
			q.events[index] = e
			q.lock.Unlock()
		}
		file.Watcher.Done()
	}
	return err
}

func checkWrite(event fsnotify.Event) bool {
	eventString := event.String()
	return strings.Contains(eventString, "CHMOD") || strings.Contains(eventString, "WRITE")
}

func ProcessQueue(processChan chan error) {
	if q.processing {
		processChan <- fmt.Errorf("process busy")
		return
	}
	q.setLock(true)
	var errs []error
	for index, v := range q.events {
		err := q.slack.PostToSlack(v)
		if err != nil {
			errs = append(errs, err)
		}
		delete(q.events, index)
	}
	q.setLock(false)
	processChan <- utils.GetErrors(errs)
}

func GetEventsCount() int {
	return len(q.events)
}

func (q *queue) setLock(running bool) {
	if running {
		q.processing = true
		q.lock.RLock()
		return
	}
	q.lock.RUnlock()
	q.processing = false
}
