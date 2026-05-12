package parser

import "sync"

type Scope int

type scopeTracker map[Scope]int

var nextScope Scope
var regMut sync.Mutex

func RegisterScope() Scope {
	regMut.Lock()
	defer regMut.Unlock()
	sc := nextScope
	nextScope++
	return sc
}

func (st scopeTracker) Enter(sc Scope) {
	st[sc]++
}

func (st scopeTracker) In(sc Scope) bool {
	return st[sc] > 0
}

func (st scopeTracker) Exit(sc Scope) {
	st[sc]--
	if st[sc] <= 0 {
		delete(st, sc)
	}
}
