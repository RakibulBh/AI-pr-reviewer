package model

type ReviewCommentRequest struct {
	Body      string `json:"body"`
	CommitID  string `json:"commit_id"`
	Path      string `json:"path"`
	StartLine int    `json:"start_line,omitempty"`
	StartSide string `json:"start_side,omitempty"`
	Line      int    `json:"line"`
	Side      string `json:"side"`
}

// PR file
type PRFile struct {
	SHA       string `json:"sha"`
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
}
