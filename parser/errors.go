package parser

import "github.com/xjslang/xjs/token"

type Range struct {
	Start token.Position `json:"start"`
	End   token.Position `json:"end"`
}

type ParserError struct {
	Message string `json:"message"`
	Range   Range  `json:"range"`
	Code    string `json:"code,omitempty"`
}
