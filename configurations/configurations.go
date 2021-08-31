package configurations

import (
	"fmt"
	"strings"

	"github.com/blainemoser/errorLogging/utils"
)

// Configs are the configs for this script
var Configs *Configure

// Configure holds configurations
type Configure struct {
	URL   string
	Files map[string]map[string]string
}

// Initialize initialises the configurations for this script
func Initialize(inputs []string) error {
	configs, err := parseInputs(inputs)
	if err != nil {
		return err
	}
	err = setConfigs(configs)
	if err != nil {
		return err
	}
	return nil
}

func setConfigs(configs map[string][]string) error {
	Configs = &Configure{}
	var errs []error
	var err error
	for i, v := range configs {
		err = setConfig(i, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utils.GetErrors(errs)
}

func args() map[string]string {
	return map[string]string{
		"url":   "url",
		"files": "files",
		"u":     "url",
		"f":     "files",
	}
}

func setConfig(prop string, setting []string) error {
	switch prop {
	case "url":
		return urlConfig(setting)
	case "files":
		return fileConfig(setting)
	default:
		return fmt.Errorf("unknown prop '%s'", prop)
	}
}

func urlConfig(setting []string) error {
	if setting[0] == "" {
		return fmt.Errorf("no url config provided")
	}
	Configs.URL = setting[0]
	return nil
}

func fileConfig(setting []string) error {
	Configs.Files = make(map[string]map[string]string)
	for _, v := range setting {
		if v == " " {
			continue
		}
		fileDetails(v)
	}
	if len(Configs.Files) < 1 {
		return fmt.Errorf("no file paths provided")
	}
	return nil
}

func fileDetails(setting string) {
	splitThis := strings.Split(setting, "@")
	fileName := splitThis[0]
	var fileFilters map[string]string
	if len(splitThis) > 1 {
		fileFilters = getFileFilters(splitThis[1])
	}
	Configs.Files[fileName] = fileFilters
}

func getFileFilters(setting string) map[string]string {
	filters := strings.Split(setting, ",")
	result := map[string]string{}
	var splitType []string
	for _, v := range filters {
		if strings.Contains(v, "#") {
			splitType = strings.Split(v, "#")
		} else {
			splitType[0] = v
			splitType[1] = "INFO"
		}
		result[splitType[0]] = splitType[1]
	}

	return result
}

func parseInputs(inputs []string) (map[string][]string, error) {
	args := args()
	result := map[string][]string{}
	var err error
	var curIndex string
	for i := 0; i < len(inputs); i++ {
		v := strings.Trim(inputs[i], " ")
		if strings.Contains(v, "=") {
			configs := strings.Split(v, "=")
			err = appendConfig(&curIndex, args, &result, configs[1], configs[0])
		} else {
			err = appendConfig(&curIndex, args, &result, v, v)
		}
		if err != nil {
			return nil, err
		}
	}
	err = checkConfigs(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func appendConfig(curIndex *string, args map[string]string, result *map[string][]string, value, index string) error {
	removeDashes(&index)
	if args[index] == "" {
		if (*result)[*curIndex] == nil {
			err := fmt.Errorf("the %s argument does not exist", index)
			return err
		}
		(*result)[*curIndex] = append((*result)[*curIndex], value)
	} else {
		*curIndex = args[index]
		(*result)[*curIndex] = make([]string, 0)
	}
	return nil
}

func removeDashes(input *string) {
	result := strings.ReplaceAll(*input, "-", "")
	*input = result
}

func checkConfigs(configs map[string][]string) error {
	var errs []error
	if configs["url"] == nil || len(configs["url"]) < 1 {
		errs = append(errs, fmt.Errorf("url config not provided"))
	}
	if len(configs["url"]) > 1 {
		errs = append(errs, fmt.Errorf("too many url arguments; expected 1, got %d", len(configs["url"])))
	}
	if configs["files"] == nil || len(configs["files"]) < 1 {
		errs = append(errs, fmt.Errorf("expected at least one file path, none provided"))
	}
	return utils.GetErrors(errs)
}
