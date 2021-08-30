package files

import (
	"fmt"
	"strings"

	"github.com/blainemoser/errorLogging/configurations"
	"github.com/blainemoser/errorLogging/filewatcher"
	"github.com/blainemoser/errorLogging/utils"
	"github.com/blainemoser/hashtable"
	"github.com/fsnotify/fsnotify"
)

type File struct {
	Watcher *filewatcher.FileWatcher
	Filters []string
}

type FileList struct {
	names    []string
	files    *hashtable.Hashtable
	notifier *fsnotify.Watcher
}

var List *FileList

func (fl *FileList) Push(name string, filters []string) error {
	watcher, err := filewatcher.Initialize(name)
	if err != nil {
		return err
	}
	file := &File{
		Watcher: watcher,
		Filters: filters,
	}
	fl.names = append(fl.names, name)
	fl.files.Put(name, file)
	return nil
}

func Initialize(notifier *fsnotify.Watcher) error {
	List = &FileList{
		names:    make([]string, 0),
		files:    hashtable.New(),
		notifier: notifier,
	}
	files := configurations.Configs.Files
	var errs []error
	var err error
	for fileName, fileFilters := range files {
		err = List.Push(fileName, fileFilters)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utils.GetErrors(errs)
}

func (fl *FileList) Walk() error {
	var errs []error
	var err error
	for _, fileName := range fl.names {
		err = fl.notifier.Add(fileName)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utils.GetErrors(errs)
}

func (fl *FileList) GetFile(path string) (*File, error) {
	fileInterface := fl.files.Get(path)
	return fileFromInterface(fileInterface)
}

func fileFromInterface(fileInterface interface{}) (*File, error) {
	if file, ok := fileInterface.(*File); ok {
		return file, nil
	}
	return nil, fmt.Errorf("file not expected type")
}

func (f *File) Filtered(e *filewatcher.Event) bool {
	if len(f.Filters) < 1 {
		return false
	}
	filtered := true
	for _, filter := range f.Filters {
		if strings.Contains(e.GetContent(), filter) {
			filtered = false
			break
		}
	}
	return filtered
}
