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

func Test() error {
	return sh.RunV("go", "test", "./...", "-timeout", "5s")
}

func TestRace() error {
	return sh.RunV("go", "test", "./...", "-race", "-timeout", "30s")
}

func Bench() error {
	return sh.RunV("go", "test", "./...", "-bench=.", "-benchtime=3s", "-run=^$")
}
