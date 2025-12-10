package sourcemap

import (
	"strings"
	"testing"
)

func TestSourceMapperBasic(t *testing.T) {
	m := New()
	m.AddMapping(1, 1)
	m.AdvanceColumn(5)
	m.AddMapping(1, 10)
	m.AdvanceColumn(3)

	sm := m.SourceMap()
	if sm.Version != 3 {
		t.Errorf("Version = %d, want 3", sm.Version)
	}
	if sm.Mappings == "" {
		t.Error("Mappings should not be empty")
	}
}

func TestSourceMapperWithNames(t *testing.T) {
	m := New()
	m.AddNamedMapping(1, 1, "foo")
	m.AdvanceColumn(3)
	m.AddNamedMapping(1, 5, "bar")
	m.AdvanceColumn(3)
	m.AddNamedMapping(2, 1, "foo") // Same name, should reuse index

	sm := m.SourceMap()
	if len(sm.Names) != 2 {
		t.Errorf("Names count = %d, want 2", len(sm.Names))
	}
	if sm.Names[0] != "foo" || sm.Names[1] != "bar" {
		t.Errorf("Names = %v, want [foo, bar]", sm.Names)
	}
}

func TestSourceMapperMultiLine(t *testing.T) {
	m := New()

	// Line 0
	m.AddMapping(1, 1)
	m.AdvanceColumn(5)

	// Line 1
	m.AdvanceLine()
	m.AddMapping(2, 1)
	m.AdvanceColumn(3)

	// Line 2
	m.AdvanceLine()
	m.AddMapping(3, 1)

	sm := m.SourceMap()
	// Mappings should contain semicolons for line breaks
	if !strings.Contains(sm.Mappings, ";") {
		t.Error("Mappings should contain semicolons for line breaks")
	}
}

func TestSourceMapperAdvanceString(t *testing.T) {
	m := New()

	m.AdvanceString("hello")
	if m.generatedColumn != 5 {
		t.Errorf("After 'hello': column = %d, want 5", m.generatedColumn)
	}

	m.AdvanceString("\n")
	if m.generatedLine != 1 || m.generatedColumn != 0 {
		t.Errorf("After newline: line=%d col=%d, want line=1 col=0",
			m.generatedLine, m.generatedColumn)
	}

	m.AdvanceString("ab\ncd\nef")
	if m.generatedLine != 3 || m.generatedColumn != 2 {
		t.Errorf("After 'ab\\ncd\\nef': line=%d col=%d, want line=3 col=2",
			m.generatedLine, m.generatedColumn)
	}
}

func TestSourceMapperAdvanceStringWindowsLineEndings(t *testing.T) {
	m := New()

	// Test single Windows line ending (\r\n)
	m.AdvanceString("hello\r\nworld")
	if m.generatedLine != 1 || m.generatedColumn != 5 {
		t.Errorf("After 'hello\\r\\nworld': line=%d col=%d, want line=1 col=5",
			m.generatedLine, m.generatedColumn)
	}

	// Test multiple Windows line endings
	m = New()
	m.AdvanceString("line1\r\nline2\r\nline3")
	if m.generatedLine != 2 || m.generatedColumn != 5 {
		t.Errorf("After multiple \\r\\n: line=%d col=%d, want line=2 col=5",
			m.generatedLine, m.generatedColumn)
	}
}

func TestSourceMapperAdvanceStringOldMacLineEndings(t *testing.T) {
	m := New()

	// Test single old Mac line ending (\r)
	m.AdvanceString("hello\rworld")
	if m.generatedLine != 1 || m.generatedColumn != 5 {
		t.Errorf("After 'hello\\rworld': line=%d col=%d, want line=1 col=5",
			m.generatedLine, m.generatedColumn)
	}

	// Test multiple old Mac line endings
	m = New()
	m.AdvanceString("line1\rline2\rline3")
	if m.generatedLine != 2 || m.generatedColumn != 5 {
		t.Errorf("After multiple \\r: line=%d col=%d, want line=2 col=5",
			m.generatedLine, m.generatedColumn)
	}
}

func TestSourceMapperAdvanceStringMixedLineEndings(t *testing.T) {
	m := New()

	// Test mixed line endings (Unix, Windows, old Mac)
	m.AdvanceString("unix\nwindows\r\noldmac\rend")
	if m.generatedLine != 3 || m.generatedColumn != 3 {
		t.Errorf("After mixed line endings: line=%d col=%d, want line=3 col=3",
			m.generatedLine, m.generatedColumn)
	}

	// Test edge case: \r at end of string followed by new string starting with \n
	m = New()
	m.AdvanceString("hello\r")
	if m.generatedLine != 1 || m.generatedColumn != 0 {
		t.Errorf("After 'hello\\r': line=%d col=%d, want line=1 col=0",
			m.generatedLine, m.generatedColumn)
	}
	// This should not treat standalone \n as part of \r\n sequence
	m.AdvanceString("\nworld")
	if m.generatedLine != 2 || m.generatedColumn != 5 {
		t.Errorf("After separate '\\nworld': line=%d col=%d, want line=2 col=5",
			m.generatedLine, m.generatedColumn)
	}
}
