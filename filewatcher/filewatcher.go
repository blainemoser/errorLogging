package filewatcher

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileWatcher struct {
	path   string
	length int64
	diff   int64
	file   *os.File
	stat   fs.FileInfo
}

type Event struct {
	start   int64 // the last byte
	end     int64 // the ending byte
	content string
	date    time.Time
	file    string
}

func Initialize(name string) (*FileWatcher, error) {
	f := &FileWatcher{}
	err := f.FirstOpen(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FileWatcher) FirstOpen(name string) error {
	path, err := filepath.EvalSymlinks(name)
	if err != nil {
		return err
	}
	f.path = path
	err = f.Start()
	if err != nil {
		return err
	}
	defer f.Done()
	return nil
}

func (f *FileWatcher) AppendedTo() bool {
	return f.diff > 0
}

func (f *FileWatcher) GetEvent() *Event {
	return &Event{
		start:   f.length,
		end:     f.length + f.diff,
		content: f.getNewContent(),
		date:    time.Now(),
		file:    f.Path(),
	}
}

func (f *FileWatcher) getNewContent() string {
	content, err := f.contentSinceLast()
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(content)
}

func (f *FileWatcher) Start() error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	f.file = file
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	f.stat = stat
	f.diff = stat.Size() - f.length
	return nil
}

func (f *FileWatcher) Path() string {
	return f.path
}

func (f *FileWatcher) Done() {
	f.length = f.stat.Size()
	f.file.Close()
}

func (f *FileWatcher) contentSinceLast() ([]byte, error) {
	f.file.Seek(f.length-1, 1)
	var content []byte
	var err error
	content, err = ioutil.ReadAll(f.file)
	if err != nil {
		return content, err
	}
	return content, nil
}

func (e *Event) GetContent() string {
	return e.content
}

func (e *Event) GetPayloadText() ([]byte, error) {
	data := map[string]interface{}{
		"text": "File: " + e.file + "\nDate: " + e.date.Format("02 January 2006 15:04:05") + "\n" + e.content,
	}
	return json.Marshal(data)
}
