module github.com/xjslang/xjs

go 1.24

require (
	github.com/magefile/mage v1.17.1
	github.com/stretchr/testify v1.11.1
	github.com/xorcare/golden v0.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// pre-estable, API experimental
retract (
    v0.1.0
    v0.2.0
    v0.3.0
    v0.4.0
    v0.5.0
    v0.6.0
    v0.7.0
    v0.8.0
    v0.8.1
    v0.9.1-alpha
    v0.9.2-alpha
    v0.9.3-alpha
    v0.10.0-beta
    v0.11.0-beta.1
    v0.11.0-beta.2
    v0.11.0-beta.3
    v0.11.0-beta.4
)
