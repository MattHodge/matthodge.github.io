package vale

import (
	"fmt"
	"os"
	"text/template"
)

func ResultToMarkdown(valeResults map[string][]Result) error {
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

	err := tm.Execute(os.Stdout, valeResults)

	if err != nil {
		return fmt.Errorf("unable to render pr comment: %v", err)
	}

	return nil
}
