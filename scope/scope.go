package scope

import "sync"

type Scope int

var scopes = make(map[Scope]int)
var nextScope Scope = 0
var regMut sync.RWMutex

func RegisterScope() Scope {
	regMut.Lock()
	defer regMut.Unlock()
	sc := nextScope
	scopes[sc] = 0
	nextScope++
	return sc
}

func Enter(sc Scope) {
	regMut.Lock()
	defer regMut.Unlock()
	scopes[sc]++
}

func In(sc Scope) bool {
	regMut.RLock()
	defer regMut.RUnlock()
	return scopes[sc] > 0
}

func Exit(sc Scope) {
	regMut.Lock()
	defer regMut.Unlock()
	if scopes[sc] > 0 {
		scopes[sc]--
	}
}
