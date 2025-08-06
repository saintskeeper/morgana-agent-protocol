# Tests Directory

This directory contains all test files and validation scripts for the Claude
Code Configuration system.

## Structure

```
tests/
├── README.md                    # This file
├── morgana/                     # Morgana Protocol tests
│   ├── test-morgana-integration.sh
│   ├── test-morgana-director-workflow.sh
│   └── test-qdirector-morgana.sh
├── integration/                 # Integration test examples
│   ├── hello.go                 # Example Go implementation
│   ├── hello_test.go           # Example Go tests
│   └── go.mod                  # Go module for examples
└── validation/                  # Migration and validation reports
    └── morgana-migration-validation.md
```

## Running Tests

### Morgana Protocol Tests

```bash
# Run all Morgana tests
cd ~/.claude/tests/morgana
./test-morgana-integration.sh

# Test director workflow
./test-morgana-director-workflow.sh

# Test QDIRECTOR integration
./test-qdirector-morgana.sh
```

### Integration Tests

```bash
# Test hello world example
cd ~/.claude/tests/integration
go test -v
go run hello.go
```

### Validation Reports

Check `tests/validation/` for migration reports and validation results.

## Test Requirements

- **Morgana Binary**: Built from `~/.claude/morgana-protocol`
- **Go**: For integration test examples
- **Bash**: For test execution scripts

## Adding New Tests

1. Place test scripts in appropriate subdirectory
2. Make scripts executable: `chmod +x script-name.sh`
3. Update this README with new test instructions
4. Follow naming convention: `test-[component]-[type].sh`
