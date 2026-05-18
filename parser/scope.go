package parser

import "sync"

type Scope int

type ScopeTracker map[Scope]int

var (
	nextScope Scope
	regMut    sync.Mutex
)

func RegisterScope() Scope {
	regMut.Lock()
	defer regMut.Unlock()
	sc := nextScope
	nextScope++
	return sc
}

func (st ScopeTracker) Enter(sc Scope) {
	st[sc]++
}

func (st ScopeTracker) In(sc Scope) bool {
	return st[sc] > 0
}

func (st ScopeTracker) Exit(sc Scope) {
	st[sc]--
	if st[sc] <= 0 {
		delete(st, sc)
	}
}
