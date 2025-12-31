package sourcemap

import (
	"strings"
)

// Mapping represents a single source map mapping entry.
// It connects a position in the generated code to a position in the original source.
type Mapping struct {
	// Generated position (in output)
	GeneratedLine   int
	GeneratedColumn int

	// Original position (in source)
	SourceLine   int
	SourceColumn int

	// Optional: index into the names array
	NameIndex int
	HasName   bool
}

// SourceMap represents a Source Map v3 structure.
type SourceMap struct {
	Version        int      `json:"version"`
	File           string   `json:"file,omitempty"`
	SourceRoot     string   `json:"sourceRoot,omitempty"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent,omitempty"`
	Names          []string `json:"names,omitempty"`
	Mappings       string   `json:"mappings"`
}

// SourceMapper accumulates mappings during code generation.
// It tracks the current position in the generated output and allows
// adding mappings back to the original source positions.
type SourceMapper struct {
	mappings  []Mapping
	names     []string
	nameIndex map[string]int

	// Current position in generated output
	generatedLine   int
	generatedColumn int
}

// New creates a new SourceMapper for the given source file.
func New() *SourceMapper {
	return &SourceMapper{
		mappings:        []Mapping{},
		names:           []string{},
		nameIndex:       make(map[string]int),
		generatedLine:   0, // 0-based
		generatedColumn: 0, // 0-based
	}
}

// AddMapping records a mapping from the current generated position
// to the given original source position (0-based line and column from token).
func (m *SourceMapper) AddMapping(sourceLine, sourceColumn int) {
	m.mappings = append(m.mappings, Mapping{
		GeneratedLine:   m.generatedLine,
		GeneratedColumn: m.generatedColumn,
		SourceLine:      sourceLine,
		SourceColumn:    sourceColumn,
	})
}

// AddNamedMapping records a mapping with an associated name (for identifiers).
func (m *SourceMapper) AddNamedMapping(sourceLine, sourceColumn int, name string) {
	nameIdx, exists := m.nameIndex[name]
	if !exists {
		nameIdx = len(m.names)
		m.names = append(m.names, name)
		m.nameIndex[name] = nameIdx
	}

	m.mappings = append(m.mappings, Mapping{
		GeneratedLine:   m.generatedLine,
		GeneratedColumn: m.generatedColumn,
		SourceLine:      sourceLine,
		SourceColumn:    sourceColumn,
		NameIndex:       nameIdx,
		HasName:         true,
	})
}

// AdvanceColumn advances the generated column position by n characters.
func (m *SourceMapper) AdvanceColumn(n int) {
	m.generatedColumn += n
}

// AdvanceString advances the position based on the content of a string,
// handling newlines appropriately.
func (m *SourceMapper) AdvanceString(s string) {
	i := 0
	for i < len(s) {
		switch s[i] {
		case '\r':
			// Handle \r\n (Windows) or \r (old Mac)
			if i+1 < len(s) && s[i+1] == '\n' {
				i++ // skip the '\n' in '\r\n'
			}
			m.generatedLine++
			m.generatedColumn = 0
		case '\n':
			m.generatedLine++
			m.generatedColumn = 0
		default:
			m.generatedColumn++
		}
		i++
	}
}

// NewLine advances to the next line in the generated output.
func (m *SourceMapper) AdvanceLine() {
	m.generatedLine++
	m.generatedColumn = 0
}

// SourceMap produces the final SourceMap structure.
func (m *SourceMapper) SourceMap() *SourceMap {
	return &SourceMap{
		Version:  3,
		Names:    m.names,
		Mappings: m.encodeMappings(),
	}
}

// encodeMappings converts all mappings to VLQ-encoded format.
func (m *SourceMapper) encodeMappings() string {
	if len(m.mappings) == 0 {
		return ""
	}

	var result strings.Builder

	// State for delta encoding (source maps use relative values)
	prevGeneratedColumn := 0
	prevSourceIndex := 0
	prevSourceLine := 0
	prevSourceColumn := 0
	prevNameIndex := 0

	currentLine := 0
	// Track how many segments have been written on the current line
	// to decide whether to insert a comma without calling result.String().
	segmentsInCurrentLine := 0

	for _, mapping := range m.mappings {
		// Add semicolons for new lines
		for currentLine < mapping.GeneratedLine {
			result.WriteByte(';')
			currentLine++
			prevGeneratedColumn = 0 // Reset column for new line
			segmentsInCurrentLine = 0
		}

		// Add comma separator if not the first segment on this line
		if segmentsInCurrentLine > 0 {
			result.WriteByte(',')
		}

		// Encode segment with delta values
		// Field 1: Generated column (delta)
		result.WriteString(encodeVLQ(mapping.GeneratedColumn - prevGeneratedColumn))
		prevGeneratedColumn = mapping.GeneratedColumn

		// Field 2: Source index (delta) - always 0 for single source
		result.WriteString(encodeVLQ(0 - prevSourceIndex))
		prevSourceIndex = 0

		// Field 3: Source line (delta)
		result.WriteString(encodeVLQ(mapping.SourceLine - prevSourceLine))
		prevSourceLine = mapping.SourceLine

		// Field 4: Source column (delta)
		result.WriteString(encodeVLQ(mapping.SourceColumn - prevSourceColumn))
		prevSourceColumn = mapping.SourceColumn

		// Field 5: Name index (delta) - optional
		if mapping.HasName {
			result.WriteString(encodeVLQ(mapping.NameIndex - prevNameIndex))
			prevNameIndex = mapping.NameIndex
		}

		// We have written a segment on this line
		segmentsInCurrentLine++
	}

	return result.String()
}
