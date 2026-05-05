package parser

type scope int

const (
	BlockScope scope = iota
)

type scopeTracker map[scope]int

func (st scopeTracker) Enter(sc scope) {
	st[sc]++
}

func (st scopeTracker) In(sc scope) bool {
	return st[sc] > 0
}

func (st scopeTracker) Exit(sc scope) {
	st[sc]--
	if st[sc] <= 0 {
		delete(st, sc)
	}
}
