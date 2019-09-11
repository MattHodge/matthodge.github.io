package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/MattHodge/matthodge.github.io/.github/actions/vale/pkg/github"
	"github.com/MattHodge/matthodge.github.io/.github/actions/vale/pkg/vale"
)

func main() {
	// Uses a default config path, but allows a user to override it from a file in their repository
	configFilePath := "/etc/vale/.vale.ini"

	userConfigFilePath := os.Getenv("INPUT_CONFIGFILEPATH")

	if fileExists(userConfigFilePath) {
		configFilePath = userConfigFilePath
	}

	lintUnchangedFiles := os.Getenv("INPUT_LINTUNCHANGEDFILES")
	lintDirectory := os.Getenv("INPUT_LINTDIRECTORY")
	fileGlob := os.Getenv("INPUT_FILEGLOB")
	githubWorkspace := os.Getenv("GITHUB_WORKSPACE")
	githubEventPath := os.Getenv("GITHUB_EVENT_PATH")
	githubEventName := os.Getenv("GITHUB_EVENT_NAME")

	fmt.Printf("lintUnchangedFiles: %s\n", lintUnchangedFiles)
	fmt.Printf("lintDirectory: %s\n", lintDirectory)
	fmt.Printf("fileGlob: %s\n", fileGlob)
	fmt.Printf("configFilePath: %s\n", configFilePath)
	fmt.Printf("githubWorkspace: %s\n", githubWorkspace)
	fmt.Printf("githubEventPath: %s\n", githubEventPath)
	fmt.Printf("githubEventName: %s\n", githubEventName)

	valeCmd := []string{
		"--no-exit",
		"--output=JSON",
		fmt.Sprintf("--config=%s", configFilePath),
		fmt.Sprintf("--glob=%s", fileGlob),
		lintDirectory,
	}

	fmt.Printf("Vale command: vale %s\n", strings.Join(valeCmd, " "))

	cmd := exec.Command("vale", valeCmd...)
	cmd.Dir = githubWorkspace
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running vale: %v", err)
		fmt.Print(string(out))
		os.Exit(1)
	}

	results := make(map[string][]vale.Result)
	err = json.Unmarshal(out, &results)
	if err != nil {
		fmt.Printf("Unable to unmarshal vale results: %v", err)
		os.Exit(1)
	}

	err = vale.ResultToMarkdown(results)

	if err != nil {
		fmt.Printf("Unable to convert results to markdown: %v", err)
		os.Exit(1)
	}

	fmt.Println("GitHub event:")
	fmt.Print(github.LoadActionsEvent(githubEventPath))
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}
