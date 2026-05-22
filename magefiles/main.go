//go:build mage

package main

import (

	// mg contains helpful utility functions, like Deps
	"fmt"

	"github.com/magefile/mage/sh"
)

func Setup() error {
	return sh.Run("git", "config", "core.hooksPath", ".githooks")
}

func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

func LintFix() error {
	return sh.RunV("golangci-lint", "run", "--fix")
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

func UpdateGoldenFiles() error {
	dirs := []string{"parser", "printer"}
	for _, dir := range dirs {
		err := sh.RunV("go", "test", fmt.Sprintf("./%s/...", dir), "-run", "^TestGoldenFiles$", "-update", "-v")
		if err != nil {
			return err
		}
	}
	return nil
}
