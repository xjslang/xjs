package scope

type Scope int

type Tracker struct {
	scopes    map[Scope]int
	nextScope Scope
}

func NewTracker() *Tracker {
	return &Tracker{
		scopes: make(map[Scope]int),
	}
}

func (t *Tracker) RegisterScope() Scope {
	sc := t.nextScope
	t.nextScope++
	return sc
}

func (t *Tracker) Enter(sc Scope) {
	t.scopes[sc]++
}

func (t *Tracker) In(sc Scope) bool {
	return t.scopes[sc] > 0
}

func (t *Tracker) Exit(sc Scope) {
	t.scopes[sc]--
	if t.scopes[sc] <= 0 {
		delete(t.scopes, sc)
	}
}
