package github

import (
	"fmt"
	"io/ioutil"
)

func LoadActionsEvent(filePath string) (string, error) {
	c, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("unable to find github event at %s", filePath)
	}

	return string(c), nil
}
