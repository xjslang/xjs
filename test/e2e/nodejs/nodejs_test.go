//go:build e2e

package nodejs_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestSourceMapsWithNodeJS(t *testing.T) {
	// Check that Node.js is available
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("Node.js not found, skipping e2e tests")
	}

	tests := []struct {
		name         string
		xjsCode      string
		expectedLine int // expected line in the original stack trace
	}{
		{
			name: "nested function call error",
			xjsCode: `function outer() {
    return inner();
}

function inner() {
    nonExistentFunction();
}

outer();`,
			expectedLine: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Parse and compile with source maps
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Build(tt.xjsCode)
			program, err := p.ParseProgram()

			if err != nil {
				t.Fatalf("Parser errors: %v", err)
			}

			result := compiler.New().WithSourceMap().Compile(program)

			// Configure source map metadata
			result.SourceMap.Sources = []string{"test.xjs"}
			result.SourceMap.SourcesContent = []string{tt.xjsCode}

			// 2. Write temporary files
			tmpDir := t.TempDir()
			jsFile := filepath.Join(tmpDir, "test.js")
			mapFile := filepath.Join(tmpDir, "test.js.map")

			// Add source map reference to the JS code
			jsCode := result.Code + "\n//# sourceMappingURL=test.js.map"
			srcMap, err := json.Marshal(result.SourceMap)
			if err != nil {
				t.Fatalf("Failed to marshal source map: %v", err)
			}

			if err := os.WriteFile(jsFile, []byte(jsCode), 0644); err != nil {
				t.Fatalf("Failed to write JS file: %v", err)
			}
			if err := os.WriteFile(mapFile, []byte(srcMap), 0644); err != nil {
				t.Fatalf("Failed to write source map file: %v", err)
			}

			// 3. Run with Node.js
			cmd := exec.Command("node", "--enable-source-maps", jsFile)
			output, err := cmd.CombinedOutput()

			// 4. Check the stack trace
			if err == nil {
				t.Fatal("Expected error but got none")
			}

			stackTrace := string(output)
			t.Logf("Stack trace:\n%s", stackTrace)

			// Check that the stack trace points to the correct line
			expectedLineRef := fmt.Sprintf(":%d:", tt.expectedLine)
			if !strings.Contains(stackTrace, expectedLineRef) {
				t.Errorf("Stack trace doesn't reference line %d. Expected to find '%s' in:\n%s",
					tt.expectedLine, expectedLineRef, stackTrace)
			}
		})
	}
}
