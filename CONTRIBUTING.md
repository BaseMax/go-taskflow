# Contributing to TaskFlow

Thank you for your interest in contributing to TaskFlow! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Git
- Make (optional, but recommended)

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-taskflow
   cd go-taskflow
   ```

3. Install dependencies:
   ```bash
   make deps
   # or
   go mod download
   ```

4. Build the project:
   ```bash
   make build
   # or
   go build -o taskflow main.go
   ```

## Development Workflow

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Testing

```bash
# Run tests
make test

# Run example workflows
./taskflow run examples/simple.yaml
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code:
  ```bash
  make fmt
  # or
  go fmt ./...
  ```

- Run linters if available:
  ```bash
  make lint
  ```

## Project Structure

```
go-taskflow/
├── cmd/taskflow/          # CLI commands
│   ├── root.go           # Root command
│   └── run.go            # Run command
├── pkg/
│   ├── types/            # Core type definitions
│   ├── parser/           # YAML parser
│   ├── executor/         # Task executors
│   └── workflow/         # Workflow engine
├── examples/             # Example workflows
├── main.go              # Entry point
├── Makefile             # Build automation
└── README.md            # Documentation
```

## Adding New Features

### Adding a New Task Type

1. Add the new task type to `pkg/types/types.go` if new fields are needed
2. Implement the executor in `pkg/executor/executor.go`
3. Add handling in the `Execute` method's switch statement
4. Create an example workflow in `examples/`
5. Update the README.md documentation

### Example: Adding a Database Task Type

```go
// In pkg/executor/executor.go

case "database":
    output, err = e.executeDatabase(ctx, task)

// Implement the executor
func (e *Executor) executeDatabase(ctx context.Context, task types.Task) (string, error) {
    // Implementation here
}
```

### Adding a New CLI Command

1. Create a new file in `cmd/taskflow/`
2. Define the command using Cobra:
   ```go
   var myCmd = &cobra.Command{
       Use:   "mycommand",
       Short: "Description",
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation
       },
   }
   
   func init() {
       rootCmd.AddCommand(myCmd)
   }
   ```

## Submitting Changes

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in present tense (Add, Fix, Update, etc.)
- Reference issue numbers if applicable

Examples:
- `Add support for PostgreSQL tasks`
- `Fix retry logic for HTTP requests`
- `Update documentation for parallel execution`

### Pull Request Process

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "Add your feature"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Open a Pull Request on GitHub

5. Ensure all checks pass

6. Wait for review and address any feedback

### Pull Request Guidelines

- Provide a clear description of the changes
- Include examples or test cases if applicable
- Update documentation as needed
- Keep changes focused and atomic
- Ensure backwards compatibility when possible

## Reporting Issues

When reporting issues, please include:

- TaskFlow version (`./taskflow --version`)
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Sample workflow file (if applicable)

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Maintain a professional tone

## Questions?

If you have questions about contributing, feel free to:
- Open an issue on GitHub
- Reach out to the maintainers

## License

By contributing to TaskFlow, you agree that your contributions will be licensed under the GPL-3.0 License.

Thank you for contributing to TaskFlow! 🎉
