package configurations

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blainemoser/errorLogging/utils"
)

func TestInitialize(t *testing.T) {
	err := Initialize([]string{
		"-u", "www.something.com/hooks/1234/5678",
		"-f",
		"/path/to/file/one@ERROR#ERROR,CRITICAL#ERROR",
		"/path/to/file/two@DEBUG#WARNING,INFO#INFO",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = checkTestedConfigs()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitializeError(t *testing.T) {
	err := Initialize([]string{
		"-f",
		"/path/to/file/one@ERROR#ERROR,CRITICAL#ERROR",
		"/path/to/file/two@DEBUG#WARNING,INFO#INFO",
	})
	if err == nil {
		t.Fatal(fmt.Errorf("expected to get error for missing url"))
		return
	}
	expectErr := "url config not provided"
	if !strings.Contains(err.Error(), expectErr) {
		t.Fatal(fmt.Errorf("expected error to be '%s', got '%s'", expectErr, err.Error()))
	}
}

func TestInitializeFullWrongConfigs(t *testing.T) {
	err := Initialize([]string{
		"-h",
		"/no/setting/h",
		"-u", "www.something.com/hooks/1234/5678",
		"-f",
		"/path/to/file/one@ERROR#ERROR,CRITICAL#ERROR",
		"/path/to/file/two@DEBUG#WARNING,INFO#INFO",
	})
	if err == nil {
		t.Fatal(fmt.Errorf("expected error for incorrect parameter"))
	}
	expectErr := "the h argument does not exist"
	if err.Error() != expectErr {
		t.Fatal(fmt.Errorf("expected error to be '%s', got '%s'", expectErr, err.Error()))
	}
}

func checkTestedConfigs() error {
	url := "www.something.com/hooks/1234/5678"
	var errs []error
	if Configs.URL != url {
		errs = append(errs, fmt.Errorf("expected url to be %s, got %s", url, Configs.URL))
	}
	checkTestedFileConfigs(&errs)
	return utils.GetErrors(errs)
}

func checkTestedFileConfigs(errs *[]error) {
	expectsFiles := getExpectsFiles()
	for fileName, settings := range expectsFiles {
		if Configs.Files[fileName] == nil {
			*errs = append(*errs, fmt.Errorf("expected %s to be in file names", fileName))
			continue
		}
		for settingProp, setting := range settings {
			if Configs.Files[fileName][settingProp] == "" {
				*errs = append(*errs, fmt.Errorf("expected %s to be in settings for file %s", settingProp, settingProp))
				continue
			}
			if setting != Configs.Files[fileName][settingProp] {
				*errs = append(*errs, fmt.Errorf("expected %s to be %s in settings for file %s, got %s", settingProp, setting, fileName, Configs.Files[fileName][settingProp]))
			}
		}
	}
}

func getExpectsFiles() map[string]map[string]string {
	return map[string]map[string]string{
		"/path/to/file/one": {
			"ERROR":    "ERROR",
			"CRITICAL": "ERROR",
		},
		"/path/to/file/two": {
			"DEBUG": "WARNING",
			"INFO":  "INFO",
		},
	}
}
