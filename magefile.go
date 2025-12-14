//go:build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test runs all project tests (unit and integration)
func Test() error {
	fmt.Println("🚀 Starting XJS Test Suite")
	fmt.Println("==========================")
	fmt.Println()
	mg.SerialDeps(TestUnit, TestE2E, TestIntegration)
	fmt.Println()
	fmt.Println("🎉 All tests completed successfully!")
	return nil
}

// TestUnit runs only unit tests (excluding integration tests)
func TestUnit() error {
	fmt.Println("🧪 Running unit tests...")
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some unit tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All unit tests passed!")
	return nil
}

// TestIntegration runs integration tests (middleware, transpilation, tolerant mode, etc.)
func TestIntegration() error {
	fmt.Println("🧪 Running integration tests...")
	err := sh.RunV("go", "test", "-v", "-tags=integration", "./test/integration/...")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some integration tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All integration tests passed!")
	return nil
}

// TestE2E runs end-to-end tests with Node.js
func TestE2E() error {
	fmt.Println("🧪 Running end-to-end tests with Node.js...")

	// Check if Node.js is available
	if !commandExists("node") {
		fmt.Println("⚠️  Node.js not found, skipping e2e tests...")
		fmt.Println("   Install Node.js to run these tests: https://nodejs.org/")
		return nil
	}

	err := sh.RunV("go", "test", "-v", "-tags=e2e", "./test/e2e/...")
	if err != nil {
		fmt.Println()
		fmt.Println("❌ Some e2e tests failed.")
		return err
	}
	fmt.Println()
	fmt.Println("✅ All e2e tests passed!")
	return nil
}

// Bench runs only benchmarks
func Bench() error {
	fmt.Println("⚡ Running benchmarks...")
	return sh.RunV("go", "test", "-bench=.", "-benchmem", "-tags=integration", "./test/integration")
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

// Docs starts a local documentation server and opens it in browser
func Docs() error {
	fmt.Println("📚 Starting documentation server...")

	// Check if pkgsite is installed
	if !commandExists("pkgsite") {
		fmt.Println("📦 Installing pkgsite...")
		if err := sh.RunV("go", "install", "golang.org/x/pkgsite/cmd/pkgsite@latest"); err != nil {
			return fmt.Errorf("failed to install pkgsite: %w", err)
		}
	}

	// Start pkgsite server in background
	fmt.Println("🌐 Starting pkgsite server on http://localhost:8080")
	cmd := exec.Command("pkgsite", "-http=localhost:8080")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start pkgsite: %w", err)
	}

	fmt.Println("💡 Press Ctrl+C to stop the documentation server")

	// Wait for the process (pkgsite will keep running until user stops it)
	return cmd.Wait()
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

// InstallHooks configures Git hooks for the project
func InstallHooks() error {
	fmt.Println("🔗 Installing Git hooks...")
	if err := sh.RunV("git", "config", "core.hooksPath", ".githooks"); err != nil {
		return fmt.Errorf("failed to set hooks path: %w", err)
	}
	if err := sh.RunV("chmod", "+x", ".githooks/pre-push"); err != nil {
		return fmt.Errorf("failed to make pre-push executable: %w", err)
	}
	fmt.Println("✅ Git hooks installed successfully!")
	return nil
}

// commandExists checks if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
