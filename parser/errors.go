package parser

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type ParserError struct {
	Message string `json:"message"`
	Range   Range  `json:"range"`
	Code    string `json:"code,omitempty"`
}
