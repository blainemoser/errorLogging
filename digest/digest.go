package digest

import (
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
	Ready      bool
	slack      *slack.Slack
	wakeUp     chan bool
}

var Q *queue

func Start() {
	Q = &queue{
		events: make(map[int]*filewatcher.Event),
		slack:  slack.NewSlack(),
		Ready:  false,
		wakeUp: make(chan bool, 1),
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
			Q.lock.Lock()
			index = len(Q.events) + 1
			Q.events[index] = e
			Q.lock.Unlock()
		}
		file.Watcher.Done()
	}

	if err != nil {
		return err
	}
	// This will trigger the ready state if the process was asleep
	// at the time of this addition
	if !Q.processing && len(Q.events) > 0 {
		Q.wakeUp <- true
	}
	return nil
}

func checkWrite(event fsnotify.Event) bool {
	eventString := event.String()
	return strings.Contains(eventString, "CHMOD") || strings.Contains(eventString, "WRITE")
}

func ProcessQueue(errChan chan error, cond *sync.Cond) {
	if !Q.processing {
		<-Q.wakeUp
	}
	Q.Ready = false
	Q.setLock(true)
	var errs []error
	for index, v := range Q.events {
		err := Q.slack.PostToSlack(v)
		if err != nil {
			errs = append(errs, err)
		}
		delete(Q.events, index)
	}
	Q.setLock(false)
	Q.Ready = true
	cond.Signal()
	errChan <- utils.GetErrors(errs)
}

func GetEventsCount() int {
	return len(Q.events)
}

func (q *queue) setLock(running bool) {
	if running {
		Q.lock.RLock()
		Q.processing = true
		return
	}
	Q.lock.RUnlock()
	Q.processing = false
}
