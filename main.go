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

var watcher *fsnotify.Watcher
var err error

func main() {
	watcher := initialise()
	defer watcher.Close()
	done := make(chan bool)
	go runWatch()
	for range time.Tick(time.Second * 1) {
		digest.ProcessQueue()
	}
	<-done
}

func runWatch() {
	for {
		err := files.List.Walk()
		if err != nil {
			log.Fatal(err)
		}
		select {
		// watch for events
		case event := <-watcher.Events:
			digest.AddToQueue(event)
		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
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

func initFileWatcher(watcher *fsnotify.Watcher) {
	err := files.Initialize(watcher)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func initialise() *fsnotify.Watcher {
	args := os.Args[1:]
	initConifgs(args)
	digest.Start()
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	initFileWatcher(watcher)
	return watcher
}
