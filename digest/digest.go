package digest

import (
	"log"
	"strings"
	"sync"

	"github.com/blainemoser/errorLogging/files"
	"github.com/blainemoser/errorLogging/filewatcher"
	"github.com/blainemoser/errorLogging/slack"
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

func AddToQueue(event fsnotify.Event) {
	var index int
	if checkWrite(event) {
		file, err := files.List.GetFile(event.Name)
		if err != nil {
			log.Println(err)
			return
		}
		file.Watcher.Start()
		e := file.Watcher.GetEvent()
		if file.Watcher.AppendedTo() && !file.Filtered(e) {
			q.lock.Lock()
			index = len(q.events) + 1
			q.events[index] = e
			q.lock.Unlock()
		}
		file.Watcher.Done()
	}
}

func checkWrite(event fsnotify.Event) bool {
	eventString := event.String()
	return strings.Contains(eventString, "CHMOD") || strings.Contains(eventString, "WRITE")
}

func ProcessQueue() {
	if q.processing {
		return
	}
	q.setLock(true)
	for index, v := range q.events {
		err := q.slack.PostToSlack(v)
		if err != nil {
			log.Println(err)
		}
		delete(q.events, index)
	}
	q.setLock(false)
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
