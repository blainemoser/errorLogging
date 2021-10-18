package filewatcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileWatcher struct {
	path   string
	length int64
	diff   int64
	file   *os.File
	stat   os.FileInfo
}

type Event struct {
	start   int64 // the last byte
	end     int64 // the ending byte
	content string
	date    time.Time
	file    string
	typeOf  string
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

func (e *Event) SetTypeOf(typeOf string) {
	e.typeOf = typeOf
}

func (e *Event) GetColor() string {
	switch strings.ToUpper(e.typeOf) {
	case "INFO":
		return "#ccf1fb"
	case "ERROR":
		return "#e2475b"
	case "CRITICAL":
		return "#e2475b"
	case "DEBUG":
		return "#ecb230"
	case "WARNING":
		return "#ecb230"
	default:
		return "#ccf1fb"
	}
}

func (e *Event) GetPayloadText() ([]byte, error) {
	data := getBlankPayload()
	e.setColor(&data)
	e.setPretext(&data)
	e.setFooter(&data)
	err := e.setFields(&data)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func (e *Event) setColor(data *string) {
	*data = strings.Replace(*data, "{{color}}", e.GetColor(), 1)
}

func (e *Event) setPretext(data *string) {
	*data = strings.Replace(*data, "{{pretext}}", "File: "+e.file, 1)
}

func (e *Event) setFooter(data *string) {
	*data = strings.Replace(*data, "{{ts}}", strconv.FormatInt(e.date.Unix(), 10), 1)
}

func (e *Event) setFields(data *string) error {
	fields, err := fields(e.content)
	if err != nil {
		return err
	}
	*data = strings.Replace(*data, "{{fields}}", fields, 1)
	return nil
}

func getBlankPayload() string {
	return `{
	"attachments": [
			{
				"mrkdwn_in": ["text"],
				"color": "{{color}}",
				"pretext": "{{pretext}}",
				"fields": {{fields}},
				"ts": {{ts}} 
			}
		]
	}`
}

func fields(content string) (string, error) {
	contentSplit := strings.Split(content, "\n")
	result := make([]interface{}, 0)
	for _, line := range contentSplit {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		result = append(result, map[string]string{
			"value": line,
		})
	}
	ret, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
