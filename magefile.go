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

// Test ejecuta todos los tests de transpilación (equivalente a test_transpilation.sh)
func Test() error {
	fmt.Println("🚀 Starting XJS Transpilation Tests")
	fmt.Println("====================================")

	// Verificar que Node.js esté disponible
	if err := checkNodeJS(); err != nil {
		return err
	}

	fmt.Println()

	// Ejecutar todos los tests de integración en secuencia
	mg.SerialDeps(TestTranspilation, TestInline, TestErrors)

	fmt.Println()
	fmt.Println("⚡ Running benchmarks...")
	if err := sh.RunV("go", "test", "-run=^$", "-bench=BenchmarkTranspilation", "-benchmem", "./test/integration"); err != nil {
		// Los benchmarks pueden fallar pero no queremos que falle todo el test
		fmt.Println("⚠️  Some benchmarks failed, but continuing...")
	}

	fmt.Println()
	fmt.Println("🎉 All tests completed successfully!")
	return nil
}

// TestTranspilation ejecuta los tests de transpilación con fixtures
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

// TestInline ejecuta los tests de transpilación inline
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

// TestErrors ejecuta los tests de manejo de errores
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

// TestUnit ejecuta solo los tests unitarios (sin tags)
func TestUnit() error {
	fmt.Println("🧪 Running unit tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// TestAll ejecuta todos los tests del proyecto (unitarios + integración)
func TestAll() error {
	fmt.Println("🧪 Running all project tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// TestMiddleware ejecuta solo los tests de middleware
func TestMiddleware() error {
	fmt.Println("⚙️ Running middleware tests...")
	return sh.RunV("go", "test", "-v", "-run", "TestMiddleware", "./test/integration")
}

// Bench ejecuta solo los benchmarks
func Bench() error {
	fmt.Println("⚡ Running benchmarks...")
	return sh.RunV("go", "test", "-bench=.", "-benchmem", "./test/integration")
}

// BenchTranspilation ejecuta solo los benchmarks de transpilación
func BenchTranspilation() error {
	fmt.Println("⚡ Running transpilation benchmarks...")
	return sh.RunV("go", "test", "-run=^$", "-bench=BenchmarkTranspilation", "-benchmem", "./test/integration")
}

// Build compila el proyecto
func Build() error {
	fmt.Println("🔨 Building XJS...")
	return sh.RunV("go", "build", "-o", "bin/xjs", ".")
}

// Clean limpia archivos generados
func Clean() error {
	fmt.Println("🧹 Cleaning generated files...")
	return sh.Rm("bin")
}

// Install instala dependencias
func Install() error {
	fmt.Println("📦 Installing dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Tidy limpia y organiza go.mod
func Tidy() error {
	fmt.Println("🔧 Tidying go.mod...")
	return sh.RunV("go", "mod", "tidy")
}

// Lint ejecuta linting (si tienes golangci-lint instalado)
func Lint() error {
	fmt.Println("🔍 Running linter...")
	if !commandExists("golangci-lint") {
		fmt.Println("⚠️  golangci-lint not found, skipping...")
		return nil
	}
	return sh.RunV("golangci-lint", "run")
}

// Dev ejecuta tests en modo watch (requiere watchexec)
func Dev() error {
	fmt.Println("🚀 Starting development mode...")
	if !commandExists("watchexec") {
		fmt.Println("ℹ️  Install watchexec for auto-testing: brew install watchexec")
		return fmt.Errorf("watchexec not found")
	}
	return sh.RunV("watchexec", "-e", "go", "-i", "bin/", "--", "mage", "test")
}

// Release prepara una release completa
func Release() error {
	fmt.Println("🚢 Preparing release...")
	mg.SerialDeps(Clean, Install, Tidy, Lint, TestAll, Build)
	fmt.Println("🎉 Release ready!")
	return nil
}

// CI ejecuta pipeline de integración continua
func CI() error {
	fmt.Println("🔄 Running CI pipeline...")
	mg.SerialDeps(Install, Lint, TestAll)
	return nil
}

// checkNodeJS verifica que Node.js esté disponible
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

// commandExists verifica si un comando existe en el PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
