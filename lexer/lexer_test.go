package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestSkipWhitespaces(t *testing.T) {
	idNames := []string{"lorem", "ipsum", "dolor"}
	l := New(strings.NewReader(fmt.Sprintf("  %s    %s %s   ", idNames[0], idNames[1], idNames[2])))
	for i := range 3 {
		tok := l.NextToken()
		if tok.Literal != idNames[i] {
			t.Errorf("Expected %s, got %s", idNames[i], tok.Literal)
		}
	}
	tok := l.NextToken()
	if tok.Literal != "" {
		t.Errorf("Expected empty string, got %s", tok.Literal)
	}
}
