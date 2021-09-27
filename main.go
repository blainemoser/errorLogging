package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/blainemoser/errorLogging/configurations"
	"github.com/blainemoser/errorLogging/digest"
	"github.com/blainemoser/errorLogging/files"
	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher := initialise()
	defer watcher.Close()
	err := files.List.Walk()
	if err != nil {
		log.Fatal(err)
		return
	}
	go runFileWatcher(watcher)
	for {
		runQueueWatcher()
	}
}

func runFileWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			err := digest.AddToQueue(event)
			if err != nil {
				log.Println(err)
			}
		case err := <-watcher.Errors:
			if err != nil {
				log.Printf("error in watcher: %s", err.Error())
			}
		}
	}
}

func runQueueWatcher() {
	errChan := make(chan error, 1)
	cond := sync.NewCond(&sync.RWMutex{})
	go digest.ProcessQueue(errChan, cond)
	cond.L.Lock()
	for !digest.Q.Ready {
		cond.Wait()
	}
	cond.L.Unlock()
	err := <-errChan
	if err != nil {
		log.Println(err.Error())
	}
}

func initConifgs(args []string) {
	err := configurations.Initialize(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func initFileWatcher(newWatcher *fsnotify.Watcher, errChan chan error) {
	errChan <- files.Initialize(newWatcher)
}

func initialise() *fsnotify.Watcher {
	args := os.Args[1:]
	initConifgs(args)
	digest.Start()
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	errChan := make(chan error, 1)
	initFileWatcher(newWatcher, errChan)
	err = <-errChan
	close(errChan)
	if err != nil {
		log.Fatal(err)
	}
	return newWatcher
}
