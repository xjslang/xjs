package parser

import (
	"encoding/json"
	"fmt"
)

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type ParserError struct {
	Message  string   `json:"message"`
	Position Position `json:"position"`
	Code     string   `json:"code,omitempty"`
}

type ParserErrors struct {
	Errors []ParserError `json:"errors"`
	Source string        `json:"source,omitempty"`
}

func (pe ParserErrors) Error() string {
	if len(pe.Errors) == 1 {
		return fmt.Sprintf("parse error at line %d, column %d: %s",
			pe.Errors[0].Position.Line,
			pe.Errors[0].Position.Column,
			pe.Errors[0].Message)
	}
	return fmt.Sprintf("parsing failed with %d errors (first: %s)",
		len(pe.Errors), pe.Errors[0].Message)
}

func (pe ParserErrors) ToJSON() ([]byte, error) {
	return json.Marshal(pe)
}
