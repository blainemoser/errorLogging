package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blainemoser/errorLogging/digest"
	"github.com/blainemoser/errorLogging/files"
	"github.com/blainemoser/errorLogging/slack"
	"github.com/fsnotify/fsnotify"
)

func TestMain(m *testing.M) {
	InitialiseTest()
	code := m.Run()
	TearDownTest()
	os.Exit(code)
}

func InitialiseTest() {
	slack.Testing = true
	err := spinTestLog()
	if err != nil {
		panic(err)
	}
}

func TearDownTest() {
	err := os.Remove("testLog.log")
	if err != nil {
		panic(err)
	}
}

func ErrorRecovery(t *testing.T) {
	recover := recover()
	if recover != nil {
		t.Error(recover)
	}
}

func spinTestLog() error {
	err := ioutil.WriteFile("testLog.log", []byte(""), 0755)
	if err != nil {
		return err
	}
	return nil
}

func TestProg(t *testing.T) {
	// this is the main loop
	os.Args = []string{"", "-u", "www.sloogle.com", "-f", "testLog.log@ERROR#ERROR,DEBUG#WARNING"}
	watcher := initialise()
	err := files.List.Walk()
	if err != nil {
		panic(err)
	}
	go runWatch(watcher)
	startFileWrite()
	processDigest()
	wChan := make(chan error, 1)
	closeWatcher(wChan, watcher)
	err = <-wChan
	close(wChan)
	if err != nil {
		log.Println(err)
	}
}

func startFileWrite() {
	writeChan := make(chan error, 1)
	writeToFile(writeChan)
	err := <-writeChan
	close(writeChan)
	if err != nil {
		panic(err)
	}
}

func processDigest() {
	closeChan := make(chan error, 1)
	digest.ProcessQueue(closeChan)
	err := <-closeChan
	close(closeChan)
	if err != nil {
		panic(err)
	}
}

func waitForWrite(written chan error) {
	timeOut := 0
	for range time.Tick(time.Millisecond * 200) {
		timeOut++
		if digest.GetEventsCount() > 0 {
			written <- nil
			break
		}
		if timeOut > 25 {
			written <- fmt.Errorf("timed out while waiting for write event")
			break
		}
	}
}

func writeToFile(writeChan chan error) {
	err := os.WriteFile("testLog.log", []byte("[ERROR]\nSomething Happened..."), 0755)
	if err != nil {
		writeChan <- err
		return
	}
	waitForWrite(writeChan)
}

func closeWatcher(closeChan chan error, w *fsnotify.Watcher) {
	err := w.Close()
	closeChan <- err
}
