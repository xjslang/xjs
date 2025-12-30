package parser

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type ParserError struct {
	Message  string   `json:"message"`
	Position Position `json:"position"`
	Length   int      `json:"length"`
	Code     string   `json:"code,omitempty"`
}
