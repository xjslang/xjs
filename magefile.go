//go:build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test runs all project tests (unit, integration, middleware)
func Test() error {
	fmt.Println("🚀 Starting XJS Test Suite")
	fmt.Println("==========================")
	fmt.Println()
	mg.SerialDeps(TestTranspilation, TestUnit, TestMiddleware, testErrors)
	fmt.Println()
	fmt.Println("🎉 All tests completed successfully!")
	return nil
}

// testTranspilation runs transpilation tests with fixtures
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

// testErrors runs error handling tests
func testErrors() error {
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

// testUnit runs only unit tests (excluding integration tests)
func TestUnit() error {
	fmt.Println("🧪 Running unit tests...")
	return sh.RunV("go", "test", "-v", "./ast", "./internal", "./lexer", "./parser", "./token")
}

// testMiddleware runs only middleware tests
func TestMiddleware() error {
	fmt.Println("⚙️ Running middleware tests...")
	return sh.RunV("go", "test", "-v", "-run", "TestMiddleware", "./test/integration")
}

// Bench runs only benchmarks
func Bench() error {
	fmt.Println("⚡ Running benchmarks...")
	return sh.RunV("go", "test", "-bench=.", "-benchmem", "./test/integration")
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

// Release prepares a complete release
func Release() error {
	fmt.Println("🚢 Preparing release...")
	mg.SerialDeps(Clean, Install, Tidy, Lint, Test)
	fmt.Println("🎉 Release ready!")
	return nil
}

// CI runs continuous integration pipeline
func CI() error {
	fmt.Println("🔄 Running CI pipeline...")
	mg.SerialDeps(Install, Lint, Test)
	return nil
}

// commandExists checks if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
