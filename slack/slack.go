package slack

import (
	"net/http"
	"strings"

	"github.com/blainemoser/errorLogging/configurations"
	"github.com/blainemoser/errorLogging/filewatcher"
)

type Slack struct {
	url string
}

func NewSlack() *Slack {
	s := &Slack{}
	s.SetURL(configurations.Configs.URL)
	return s
}

func (s *Slack) SetURL(url string) {
	s.url = url
}

func (s *Slack) PostToSlack(e *filewatcher.Event) error {

	data, err := e.GetPayloadText()
	if err != nil {
		return err
	}

	_, err = http.Post(s.url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	return nil
}
