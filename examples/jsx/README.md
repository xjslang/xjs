**[XJS](https://github.com/xjslang/xjs) is a parsing tool that allows us to extend JavaScript with our favorite features.** And this basic example illustrates how we can extend JS to use HTML tags in our expressions. For example:

```js
let p = <p>
  "Hello, " | <strong>"Word!"</strong>
</p>

// is transpiled to
let p = (function () {
  const elem = document.createElement('p');
  elem.append(
    (function () {
      const elem = document.createDocumentFragment();
      elem.append("Hello, ");
      elem.append(
        (function () {
          const elem = document.createElement('strong');
          elem.append("World!");
          return elem;
        })(),
      );
      return elem;
    })(),
  );
  return elem;
})();
```

## How to use it

You will find more examples in [./jsx_test.go](./jsx_test.go).

```go
package main

import (
	"fmt"

	"github.com/xjslang/xjs/examples/jsx"
)

func main() {
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
	// Output: let p = (function(){const elem = ...})();
}
```
