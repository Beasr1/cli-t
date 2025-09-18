# CLI-T Test Suite

This directory contains the test infrastructure for CLI-T.

## Structure

```
test/
├── helpers/          # Test utilities and helpers
│   └── command.go    # Command execution helpers
├── testdata/         # Test fixtures and data
│   └── config.yaml   # Test configuration
├── integration/      # Integration tests (coming soon)
└── e2e/             # End-to-end tests (coming soon)
```

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only (Fast)
```bash
make test-unit
```

### With Coverage
```bash
make coverage
```

### Specific Package
```bash
make test-pkg PKG=internal/command
```

### Benchmarks
```bash
make test-bench
```

### Watch Mode
```bash
make watch
```

## Writing Tests

### Unit Test Example
```go
func TestMyFunction(t *testing.T) {
    // Arrange
    cmd := NewCommand()
    args := helpers.TestArgs("input.txt")
    
    // Act
    err := cmd.Execute(context.Background(), args)
    
    // Assert
    assert.NoError(t, err)
    assert.Contains(t, args.Stdout.(*bytes.Buffer).String(), "expected output")
}
```

### Table-Driven Test Example
```go
func TestWordCount(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        flags    []string
        expected string
        wantErr  bool
    }{
        {
            name:     "count lines",
            input:    "line1\nline2\n",
            flags:    []string{"-l"},
            expected: "2",
        },
        // more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Using Test Helpers
```go
// Create test arguments
args := helpers.TestArgs("file1.txt", "--verbose")

// Execute command and capture output
stdout, stderr, err := helpers.ExecuteCommand(t, cmd, "arg1", "arg2")
```

## Test Guidelines

1. **Isolation**: Each test should be independent
2. **Clarity**: Test names should describe what they test
3. **Coverage**: Aim for >70% coverage
4. **Speed**: Mark slow tests with `t.Skip()` in short mode
5. **Determinism**: No flaky tests - use fixed time, randomness seeds

## Mocking

For commands that need external dependencies:

```go
type MockReader struct {
    data []byte
    err  error
}

func (m *MockReader) Read(p []byte) (n int, err error) {
    if m.err != nil {
        return 0, m.err
    }
    return copy(p, m.data), io.EOF
}
```

## Integration Tests

Integration tests go in `test/integration/` and test multiple components together:

```go
// +build integration

func TestFullCommandExecution(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    // test implementation
}
```

## Benchmarks

Write benchmarks for performance-critical code:

```go
func BenchmarkLargeFileProcessing(b *testing.B) {
    data := generateTestData(1024 * 1024) // 1MB
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

## CI/CD Integration

Tests run automatically on:
- Every push to main
- Every pull request
- Can be run manually with `make ci-test`

## Debugging Tests

### Verbose Output
```bash
go test -v ./...
```

### Run Specific Test
```bash
go test -v -run TestRegistry ./internal/command
```

### Debug with Delve
```bash
dlv test ./internal/command
```

## Test Data

Place test fixtures in `testdata/` directories within each package:

```
internal/tools/wc/testdata/
├── empty.txt
├── single-line.txt
└── multi-line.txt
```

These are automatically ignored by Go's build system.