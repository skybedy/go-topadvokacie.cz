package model

type Result struct {
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Sections []Section `json:"sections"`
	Warnings []string  `json:"warnings"`
	Raw      string    `json:"raw,omitempty"`
}

type Section struct {
	Title string   `json:"title"`
	Items []string `json:"items"`
}

type Example struct {
	ID      string
	Title   string
	Summary string
	Content string
}
