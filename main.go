package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blainemoser/errorLogging/configurations"
	"github.com/blainemoser/errorLogging/digest"
	"github.com/blainemoser/errorLogging/files"
	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher := initialise()
	defer watcher.Close()
	done := make(chan bool)
	err := files.List.Walk()
	if err != nil {
		log.Fatal(err)
		return
	}
	go runWatch(watcher)
	for range time.Tick(time.Second * 1) {
		processChan := make(chan error, 1)
		digest.ProcessQueue(processChan)
		processError := <-processChan
		close(processChan)
		if processError != nil {
			log.Println(processError.Error())
		}
	}
	<-done
}

func runWatch(watcher *fsnotify.Watcher) {
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
