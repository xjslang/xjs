package parser

type scope int

const (
	blockScope scope = iota
)

type scopeTracker map[scope]int

func (st scopeTracker) enter(sc scope) {
	st[sc]++
}

func (st scopeTracker) in(sc scope) bool {
	return st[sc] > 0
}

func (st scopeTracker) exit(sc scope) {
	st[sc]--
	if st[sc] <= 0 {
		delete(st, sc)
	}
}
