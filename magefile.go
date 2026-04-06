//go:build mage

package main

import (

	// mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

func Setup() error {
	return sh.Run("git", "config", "core.hooksPath", ".githooks")
}

func Lint() error {
	return sh.RunV("golangci-lint", "run")
}
