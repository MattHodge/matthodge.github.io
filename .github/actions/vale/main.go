package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
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

	fmt.Printf("lintUnchangedFiles: %s\n", lintUnchangedFiles)
	fmt.Printf("lintDirectory: %s\n", lintDirectory)
	fmt.Printf("fileGlob: %s\n", fileGlob)
	fmt.Printf("configFilePath: %s\n", configFilePath)
	fmt.Printf("githubWorkspace: %s\n", githubWorkspace)
	fmt.Printf("githubEventPath: %s\n", githubEventPath)

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

	results := make(map[string][]ValeResult)
	err = json.Unmarshal(out, &results)
	if err != nil {
		fmt.Printf("Unable to unmarshal vale results: %v", err)
		os.Exit(1)
	}

	const prTemplate = `
{{ range $key, $value := . }}
### ` + "`{{ $key }}`" +
		`
Check Name | Line | Message | Severity
--- | --- | --- | ---
{{ range $value := . -}}
{{ .Check }} | {{ .Line }} | {{ .Message }} | {{ .Severity }}
{{ end -}}
{{ end -}}
`

	// Create a new template and parse the letter into it.
	tm := template.Must(template.New("pr-comment").Parse(prTemplate))

	err = tm.Execute(os.Stdout, results)
	if err != nil {
		fmt.Printf("Unable to render pr comment: %v", err)
		os.Exit(1)
	}

	fmt.Println("GitHub event:")
	fmt.Print(loadGitHubEvent(githubEventPath))
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadGitHubEvent(filePath string) (string, error) {
	c, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("unable to find github event at %s", filePath)
	}

	return string(c), nil
}

type ValeResult struct {
	Check       string `json:"Check"`
	Description string `json:"Description"`
	Line        int    `json:"Line"`
	Link        string `json:"Link"`
	Message     string `json:"Message"`
	Severity    string `json:"Severity"`
	Span        []int  `json:"Span"`
	Hide        bool   `json:"Hide"`
	Match       string `json:"Match"`
}
