//go:build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when no target is specified
var Default = Test

// Test runs all transpilation tests (equivalent to test_transpilation.sh)
func Test() error {
	fmt.Println("🚀 Starting XJS Transpilation Tests")
	fmt.Println("====================================")
	if err := checkNodeJS(); err != nil {
		return err
	}
	fmt.Println()
	mg.SerialDeps(TestTranspilation, TestInline, TestErrors)
	fmt.Println()
	fmt.Println("⚡ Running benchmarks...")
	if err := sh.RunV("go", "test", "-run=^$", "-bench=BenchmarkTranspilation", "-benchmem", "./test/integration"); err != nil {
		fmt.Println("⚠️  Some benchmarks failed, but continuing...")
	}
	fmt.Println()
	fmt.Println("🎉 All tests completed successfully!")
	return nil
}

// TestTranspilation runs transpilation tests with fixtures
func TestTranspilation() error {
	fmt.Println("🧪 Running transpilation tests...")
	err := sh.RunV("go", "test", "-v", "-run", "TestTranspilation$", "./test/integration")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some transpilation tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All transpilation tests passed!")
	return nil
}

// TestInline runs inline transpilation tests
func TestInline() error {
	fmt.Println("🎯 Running inline transpilation tests...")
	err := sh.RunV("go", "test", "-v", "-run", "TestTranspilationBasicInline", "./test/integration")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some inline transpilation tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All inline transpilation tests passed!")
	return nil
}

// TestErrors runs error handling tests
func TestErrors() error {
	fmt.Println("🔥 Running error handling tests...")
	err := sh.RunV("go", "test", "-v", "-run", "TestTranspilationErrors", "./test/integration")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some error handling tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All error handling tests passed!")
	return nil
}

// TestUnit runs only unit tests (without tags)
func TestUnit() error {
	fmt.Println("🧪 Running unit tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// TestAll runs all project tests (unit + integration)
func TestAll() error {
	fmt.Println("🧪 Running all project tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// TestMiddleware runs only middleware tests
func TestMiddleware() error {
	fmt.Println("⚙️ Running middleware tests...")
	return sh.RunV("go", "test", "-v", "-run", "TestMiddleware", "./test/integration")
}

// Bench runs only benchmarks
func Bench() error {
	fmt.Println("⚡ Running benchmarks...")
	return sh.RunV("go", "test", "-bench=.", "-benchmem", "./test/integration")
}

// BenchTranspilation runs only transpilation benchmarks
func BenchTranspilation() error {
	fmt.Println("⚡ Running transpilation benchmarks...")
	return sh.RunV("go", "test", "-run=^$", "-bench=BenchmarkTranspilation", "-benchmem", "./test/integration")
}

// Clean removes temporary files and cache
func Clean() error {
	fmt.Println("🧹 Cleaning temporary files and cache...")
	if err := sh.RunV("go", "clean", "-testcache"); err != nil {
		fmt.Println("Note: failed to clean test cache, continuing...")
	}
	if err := sh.RunV("go", "clean", "-modcache"); err != nil {
		fmt.Println("Note: failed to clean mod cache, continuing...")
	}
	fmt.Println("✅ Cleanup completed!")
	return nil
}

// Install installs dependencies
func Install() error {
	fmt.Println("📦 Installing dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Tidy cleans and organizes go.mod
func Tidy() error {
	fmt.Println("🔧 Tidying go.mod...")
	return sh.RunV("go", "mod", "tidy")
}

// Lint runs linting (if golangci-lint is installed)
func Lint() error {
	fmt.Println("🔍 Running linter...")
	if !commandExists("golangci-lint") {
		fmt.Println("⚠️  golangci-lint not found, skipping...")
		return nil
	}
	return sh.RunV("golangci-lint", "run")
}

// Dev runs tests in watch mode (requires watchexec)
func Dev() error {
	fmt.Println("🚀 Starting development mode...")
	if !commandExists("watchexec") {
		fmt.Println("ℹ️  Install watchexec for auto-testing: brew install watchexec")
		return fmt.Errorf("watchexec not found")
	}
	return sh.RunV("watchexec", "-e", "go", "-i", "bin/", "--", "mage", "test")
}

// Release prepares a complete release
func Release() error {
	fmt.Println("🚢 Preparing release...")
	mg.SerialDeps(Clean, Install, Tidy, Lint, TestAll)
	fmt.Println("🎉 Release ready!")
	return nil
}

// CI runs continuous integration pipeline
func CI() error {
	fmt.Println("🔄 Running CI pipeline...")
	mg.SerialDeps(Install, Lint, TestAll)
	return nil
}

// checkNodeJS verifies that Node.js is available
func checkNodeJS() error {
	if !commandExists("node") {
		fmt.Println("❌ Error: Node.js is required but not installed.")
		fmt.Println("   Please install Node.js to run transpilation tests.")
		return fmt.Errorf("node.js not found")
	}
	version, err := sh.Output("node", "--version")
	if err != nil {
		return err
	}
	fmt.Printf("✅ Node.js found: %s\n", version)
	return nil
}

// commandExists checks if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
