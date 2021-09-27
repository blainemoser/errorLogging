package slack

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jsonextract "github.com/blainemoser/JsonExtract"
	"github.com/blainemoser/errorLogging/configurations"
	"github.com/blainemoser/errorLogging/filewatcher"
	"github.com/blainemoser/errorLogging/utils"
)

var Testing bool

var TestingPosts int

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

	// We won't be sending anything in a test context
	if !Testing {
		_, err = http.Post(s.url, "application/json", strings.NewReader(string(data)))
		if err != nil {
			return err
		}
	} else {
		return checkData(data, e)
	}

	return nil
}

func checkData(data []byte, e *filewatcher.Event) error {
	var errs []error
	extract := &jsonextract.JSONExtract{
		RawJSON: string(data),
	}
	attachmentInterface, err := extract.Extract("attachments/[0]")
	if err != nil {
		return err
	}
	if attachment, ok := attachmentInterface.(map[string]interface{}); ok {
		checkPretext(attachment, &errs)
		checkFields(extract, attachment, e, &errs)
		return utils.GetErrors(errs)
	}
	return fmt.Errorf("expected attachment node in the slack payload, none found")
}

func checkPretext(attachment map[string]interface{}, errs *[]error) {
	pretext, err := utils.StringMapInterface(attachment, "pretext")
	if err != nil {
		*errs = append(*errs, err)
	} else {
		expected := "File: testLogOne.log"
		expectedTwo := "File: testLogTwo.log"
		if pretext != expected && pretext != expectedTwo {
			*errs = append(*errs, fmt.Errorf("expected pretext to be '%s' or '%s', got '%s", expected, expectedTwo, pretext))
		}
	}
}

func checkFields(extract *jsonextract.JSONExtract, attachment map[string]interface{}, e *filewatcher.Event, errs *[]error) {
	fields, err := utils.ArrayInterface(attachment, "fields")
	if err != nil {
		*errs = append(*errs, err)
	} else {
		TestingPosts += len(fields)
		for i := range fields {
			f, err := extract.Extract("attachments/[0]/fields/[" + strconv.Itoa(i) + "]/value")
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			value, err := utils.StringInterface(f)
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			checkFieldValue(value, e, errs)
		}
	}
}

func checkFieldValue(value string, e *filewatcher.Event, errs *[]error) {
	if !strings.Contains(e.GetContent(), value) {
		*errs = append(*errs, fmt.Errorf("unexpected '%s' in payload content", value))
	}
}
