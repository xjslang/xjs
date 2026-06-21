package jsx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/examples/jsx"
	"github.com/xjslang/xjs/printer"
)

func ExampleCompile() {
	// transform the input to AST
	input := `let p = <p>
		"Hello, " |
		<strong>"World!"</strong>
	</p>`
	result, err := jsx.Parse([]byte(input))
	if err != nil {
		panic(err)
	}

	// transform the AST to valid JS code
	jsCode, err := jsx.Compile(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsCode)
	// Output: let p = (function(){const elem = document.createElement('p');elem.append((function(){const elem = document.createDocumentFragment();elem.append("Hello, ");elem.append((function(){const elem = document.createElement('strong');elem.append("World!");return elem})());return elem})());return elem})();
}

func ExampleFormat() {
	// transform the input to AST
	input := `let p = <p>"Hello, " | <strong>"World!"</strong></p>`
	result, err := jsx.Parse([]byte(input))
	if err != nil {
		panic(err)
	}

	// transform the AST to properly formatted XJS code
	xjsCode, err := jsx.Format(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(xjsCode)
	// Output:
	// let p = <p>
	//   "Hello, " | <strong>
	//     "World!"
	//   </strong>
	// </p>;
}

func TestParse(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		input := "<div>'Hello, World!'</p>"
		_, err := jsx.Parse([]byte(input))
		require.Error(t, err)
		require.Equal(t, err.Error(), "[line:0, col:22] expected closing tag </div>")
	})
}

func TestCompile(t *testing.T) {
	t.Run("empty tag", func(t *testing.T) {
		input := `let p = <p></p>`
		result, err := jsx.Parse([]byte(input))
		require.NoError(t, err)
		out, err := jsx.Compile(result)
		require.NoError(t, err)
		require.Equal(t, "let p = (function(){const elem = document.createElement('p');return elem})();", out)
	})
}

func TestFormat(t *testing.T) {
	t.Run("empty tags", func(t *testing.T) {
		input := `let p = <p></p>`
		result, err := jsx.Parse([]byte(input))
		require.NoError(t, err)
		out, err := jsx.Format(result)
		require.NoError(t, err)
		require.Equal(t, "let p = <p>\n</p>;", out)
	})
	t.Run("with comments", func(t *testing.T) {
		input := `let p = <p>
		// c1
		"Hello, " |
		<strong>
		/* c2 */
		"World!"</strong></p>`
		result, err := jsx.Parse([]byte(input))
		require.NoError(t, err)

		// transform the AST to properly formatted code
		code, err := jsx.Format(result, printer.WithIndent("\t"))
		require.NoError(t, err)
		expectedCode := "let p = <p>\n\t// c1\n\t\"Hello, \" | <strong>\n\t\t/* c2 */\n\t\t\"World!\"\n\t</strong>\n</p>;"
		require.Equal(t, expectedCode, code)
	})
}
