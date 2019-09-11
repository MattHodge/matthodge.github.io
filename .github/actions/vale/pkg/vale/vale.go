package vale

type Result struct {
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
