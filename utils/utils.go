package utils

import (
	"fmt"
	"strings"
)

func GetErrors(errs []error) error {
	var errStrings []string
	if len(errs) > 0 {
		for _, e := range errs {
			errStrings = append(errStrings, e.Error())
		}
		return fmt.Errorf(strings.Join(errStrings, "; "))
	}
	return nil
}

func StringMapInterface(input map[string]interface{}, key string) (string, error) {
	if input[key] == nil {
		return "", fmt.Errorf("'%s' not found", key)
	}
	return StringInterface(input[key])
}

func StringInterface(input interface{}) (string, error) {
	if shouldBe, ok := input.(string); ok {
		return shouldBe, nil
	}
	return "", fmt.Errorf("input is not a string")
}

func ArrayInterface(input map[string]interface{}, key string) ([]interface{}, error) {
	if input[key] == nil {
		return nil, fmt.Errorf("'%s' not found", key)
	}
	if shouldBe, ok := input[key].([]interface{}); ok {
		return shouldBe, nil
	}
	return nil, fmt.Errorf("input is not an array")
}
