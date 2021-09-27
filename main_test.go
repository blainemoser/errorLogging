package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

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
		os.Remove("testLogOne.log")
		os.Remove("testLogTwo.log")
		panic(err)
	}
}

func TearDownTest() {
	err := os.Remove("testLogOne.log")
	if err != nil {
		log.Println(err)
	}
	err = os.Remove("testLogTwo.log")
	if err != nil {
		log.Println(err)
	}
}

func ErrorRecovery(t *testing.T) {
	recover := recover()
	if recover != nil {
		TearDownTest()
		t.Error(recover)
	}
}

func spinTestLog() error {
	err := ioutil.WriteFile("testLogOne.log", []byte(""), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("testLogTwo.log", []byte(""), 0755)
	if err != nil {
		return err
	}
	return nil
}

func TestProg(t *testing.T) {
	defer ErrorRecovery(t)
	// this is the main loop
	os.Args = []string{
		"",
		"-u",
		"https://hooks.slack.com/services/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",
		"-f",
		"testLogOne.log@ERROR#ERROR,DEBUG#WARNING",
		"testLogTwo.log@ERROR#ERROR,DEBUG#WARNING",
	}
	watcher := initialise()
	err := files.List.Walk()
	if err != nil {
		panic(err)
	}
	go runFileWatcher(watcher)
	go func() {
		for {
			runQueueWatcher()
		}
	}()
	err = startFileWrite()
	if err != nil {
		panic(err)
	}
	wChan := make(chan error, 1)
	closeWatcher(wChan, watcher)
	err = <-wChan
	close(wChan)
	if err != nil {
		panic(err)
	}
}

func startFileWrite() error {
	writeChan := make(chan error, 1)
	writeToFiles(writeChan)
	err := <-writeChan
	close(writeChan)
	if err != nil {
		return err
	}
	timeout := 0
	for range time.Tick(time.Millisecond * 200) {
		timeout++
		if slack.TestingPosts == 25 {
			return nil
		}
		if timeout >= 50 {
			break
		}
	}
	return fmt.Errorf("expected 25 posts to have been made before timeout, got %d", slack.TestingPosts)
}

func writeToFiles(writeChan chan error) {
	rand.Seed(time.Now().Unix())
	var randType int
	logOne, err := os.OpenFile("testLogOne.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		writeChan <- err
		return
	}
	logTwo, err := os.OpenFile("testLogTwo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		writeChan <- err
		return
	}
	defer logOne.Close()
	defer logTwo.Close()
	for i := 0; i < 25; i++ {
		randType = rand.Intn(2)
		if randType == 1 {
			a := strconv.AppendInt([]byte("[ERROR] Something Happened to Log One... at "), (time.Now().UnixNano()), 10)
			if _, err := logOne.WriteString(string(a) + "\n"); err != nil {
				writeChan <- err
				return
			}
		} else {
			a := strconv.AppendInt([]byte("[ERROR] Something Happened to Log Two... at "), (time.Now().UnixNano()), 10)
			if _, err := logTwo.WriteString(string(a) + "\n"); err != nil {
				writeChan <- err
				return
			}
		}
	}
	writeChan <- nil
}

func closeWatcher(closeChan chan error, w *fsnotify.Watcher) {
	err := w.Close()
	closeChan <- err
}
