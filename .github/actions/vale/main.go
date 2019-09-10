package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	fmt.Printf("lintUnchangedFiles: %s\n", lintUnchangedFiles)
	fmt.Printf("lintDirectory: %s\n", lintDirectory)
	fmt.Printf("fileGlob: %s\n", fileGlob)
	fmt.Printf("configFilePath: %s\n", configFilePath)
	fmt.Printf("githubWorkspace: %s\n", githubWorkspace)

	valeCmd := []string{
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
		fmt.Print(string(out))
		fmt.Printf("Error running vale: %v", err)
		os.Exit(1)
	}

	fmt.Print(string(out))
}

func fileExists(filePath string) bool {
	if _, err := os.Stat("filePath"); os.IsNotExist(err) {
		return false
	}

	return true
}
