package source

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}
