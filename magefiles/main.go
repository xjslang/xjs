//go:build mage

package main

import (
	"os"
	"os/exec"
	"strings"

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
	return sh.RunV("go", "test", "./printer", "./js", "./jsextended", "-update")
}

func CleanTestCache() error {
	return sh.RunV("go", "clean", "-testcache")
}

func Docs() error {
	return sh.RunV("pkgsite", "-http", "localhost:8081")
}

func BenchCompare(targetBranch, benchmark string) error {
	currentBranch, err := sh.Output("git", "branch", "--show-current")
	if err != nil {
		return err
	}
	currentBranch = strings.TrimSpace(currentBranch)

	// Ensure we return to the original branch even if something fails.
	defer func() {
		_ = sh.Run("git", "checkout", currentBranch)
	}()

	if err := sh.Run("git", "checkout", targetBranch); err != nil {
		return err
	}

	if err := runBenchmarkToFile(benchmark, "before.out"); err != nil {
		return err
	}

	if err := sh.Run("git", "checkout", currentBranch); err != nil {
		return err
	}

	if err := runBenchmarkToFile(benchmark, "after.out"); err != nil {
		return err
	}

	return sh.RunV("benchstat", "before.out", "after.out")
}

func runBenchmarkToFile(benchmark, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(
		"go", "test",
		"-run=^"+benchmark+"$",
		"-bench=^"+benchmark+"$",
		"-benchmem",
		"-count=10",
		"./...",
	)
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
