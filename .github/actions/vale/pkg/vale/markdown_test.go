package vale

import (
	"testing"
)

func TestResultToMarkdown(t *testing.T) {

	results := make(map[string][]Result)

	results["readme.md"] = []Result{
		{
			Check:    "write-good.Weasel",
			Line:     2,
			Message:  "'fairly' is a weasel word!",
			Severity: "warning",
		},
		{
			Check:    "write-good.Weasel",
			Line:     3,
			Message:  "'fairly' is a weasel word!",
			Severity: "warning",
		},
	}

	results["other.md"] = []Result{
		{
			Check:    "write-good.Weasel",
			Line:     2,
			Message:  "'fairly' is a weasel word!",
			Severity: "warning",
		},
		{
			Check:    "write-good.Weasel",
			Line:     3,
			Message:  "'fairly' is a weasel word!",
			Severity: "warning",
		},
	}

	err := ResultToMarkdown(results)

	if err != nil {
		t.Errorf("Unable to convert result to Markdown %v", err)
	}
}
