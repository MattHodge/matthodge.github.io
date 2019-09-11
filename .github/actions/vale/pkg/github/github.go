package github

import (
	"fmt"
	"io/ioutil"

	"github.com/google/go-github/github"
)

func LoadActionsEvent(messageType, filePath string) (interface{}, error) {
	c, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("unable to find github event at %s", filePath)
	}

	evt, err := github.ParseWebHook(messageType, c)

	if err != nil {
		return nil, fmt.Errorf("unable to parse github event: %v", err)
	}

	return evt, nil
}
