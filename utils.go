package goswish

import (
	"regexp"
	"strings"
)

func GetInstructionID(location string) (string, error) {
	regex, err := regexp.Compile(`\/([^\/]+)$`)
	if err != nil {
		return "", err
	}
	matches := regex.FindAllString(location, -1)
	if len(matches) == 0 {
		return "", nil
	}
	return strings.ReplaceAll(matches[0], "/", ""), nil
}
