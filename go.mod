module github.com/xjslang/xjs

go 1.26

require (
	github.com/dop251/goja v0.0.0-20260826172338-03e15ec872a2
	github.com/magefile/mage v1.17.1
	github.com/stretchr/testify v1.11.1
	github.com/tdewolff/parse/v2 v2.8.16
	github.com/xorcare/golden v0.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.3.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// pre-estable, API experimental
retract (
	v0.11.0-beta.4
	v0.11.0-beta.3
	v0.11.0-beta.2
	v0.11.0-beta.1
	v0.10.0-beta
	v0.9.3-alpha
	v0.9.2-alpha
	v0.9.1-alpha
	v0.8.1
	v0.8.0
	v0.7.0
	v0.6.0
	v0.5.0
	v0.4.0
	v0.3.0
	v0.2.0
	v0.1.0
)
