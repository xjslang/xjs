#!/bin/bash

# XJS Transpilation Test Runner
# This script runs transpilation tests for the XJS compiler

echo "🚀 Starting XJS Transpilation Tests"
echo "===================================="

# Check if Node.js is available
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is required but not installed."
    echo "   Please install Node.js to run transpilation tests."
    exit 1
fi

echo "✅ Node.js found: $(node --version)"
echo ""

# Run the transpilation tests
echo "🧪 Running transpilation tests..."
go test -v -run TestTranspilation

# Check if tests passed
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All transpilation tests passed!"
else
    echo ""
    echo "❌ Some transpilation tests failed."
    exit 1
fi

echo ""
echo "🎯 Running inline transpilation tests..."
go test -v -run TestTranspilationBasicInline

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All inline transpilation tests passed!"
else
    echo ""
    echo "❌ Some inline transpilation tests failed."
    exit 1
fi

echo ""
echo "🔥 Running error handling tests..."
go test -v -run TestTranspilationErrors

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All error handling tests passed!"
else
    echo ""
    echo "❌ Some error handling tests failed."
    exit 1
fi

echo ""
echo "⚡ Running benchmarks..."
go test -bench=BenchmarkTranspilation -benchmem

echo ""
echo "🎉 All tests completed successfully!"
