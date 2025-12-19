# go-taskflow

A declarative task automation tool designed for developers and ops teams.

## Overview

TaskFlow is a powerful Go-based automation engine that executes workflows defined in YAML files. It provides a clean CLI interface and supports:

- 🔧 **Multiple Task Types**: Shell commands, HTTP requests, and file operations
- 🔗 **Task Dependencies**: Define execution order with dependency chains
- ⚡ **Parallel Execution**: Run independent tasks concurrently
- 🔄 **Retry Logic**: Automatic retry with configurable delays
- 🎯 **Conditional Execution**: Skip tasks based on runtime conditions
- ⏱️ **Timeouts**: Set execution time limits for tasks
- 📝 **Variable Substitution**: Use dynamic variables across workflows
- 🎨 **Beautiful CLI Output**: Clear progress and results display

## Installation

### Build from source

```bash
git clone https://github.com/BaseMax/go-taskflow
cd go-taskflow
go build -o taskflow main.go
```

### Quick start

```bash
# Run a workflow
./taskflow run examples/simple.yaml

# Show help
./taskflow --help

# Check version
./taskflow --version
```

## Workflow Syntax

### Basic Structure

```yaml
name: My Workflow
description: Description of what this workflow does

variables:
  my_var: "value"
  another_var: "hello"

tasks:
  - name: task-name
    type: shell  # or http, file
    command: echo "Hello ${my_var}"
```

### Task Types

#### Shell Commands

```yaml
- name: run-script
  type: shell
  command: echo "Hello World" && ls -la
```

#### HTTP Requests

```yaml
- name: api-call
  type: http
  url: https://api.example.com/endpoint
  method: POST
  headers:
    Content-Type: application/json
  body: '{"key": "value"}'
```

#### File Operations

```yaml
# Write file
- name: write-config
  type: file
  file_action: write
  file_path: /tmp/config.txt
  file_content: "Configuration data"

# Read file
- name: read-config
  type: file
  file_action: read
  file_path: /tmp/config.txt

# Copy file
- name: backup-config
  type: file
  file_action: copy
  source_path: /tmp/config.txt
  dest_path: /tmp/config.backup.txt

# Delete file
- name: cleanup
  type: file
  file_action: delete
  file_path: /tmp/config.txt
```

### Advanced Features

#### Dependencies

```yaml
tasks:
  - name: setup
    type: shell
    command: echo "Setting up..."

  - name: build
    type: shell
    command: echo "Building..."
    depends_on:
      - setup

  - name: test
    type: shell
    command: echo "Testing..."
    depends_on:
      - build
```

#### Parallel Execution

```yaml
tasks:
  - name: worker-1
    type: shell
    command: echo "Worker 1"
    parallel: true

  - name: worker-2
    type: shell
    command: echo "Worker 2"
    parallel: true

  - name: worker-3
    type: shell
    command: echo "Worker 3"
    parallel: true
```

#### Retry Logic

```yaml
- name: flaky-task
  type: shell
  command: ./might-fail.sh
  retry:
    max_attempts: 3
    delay: 2s
  continue_on_error: true
```

#### Conditional Execution

```yaml
variables:
  environment: "production"

tasks:
  - name: production-only
    type: shell
    command: echo "Running in production"
    condition: ${environment} == "production"
```

#### Timeouts

```yaml
- name: long-running
  type: shell
  command: ./long-process.sh
  timeout: 5m
```

## Examples

The `examples/` directory contains several workflow examples:

- **simple.yaml** - Basic shell commands and variable usage
- **parallel.yaml** - Parallel task execution demonstration
- **http.yaml** - HTTP API calls
- **file-ops.yaml** - File operations (read, write, copy, delete)
- **retry.yaml** - Retry logic and error handling
- **ci-pipeline.yaml** - Complete CI/CD pipeline example

Run any example:

```bash
./taskflow run examples/ci-pipeline.yaml
```

## Architecture

```
go-taskflow/
├── cmd/taskflow/          # CLI commands
│   ├── root.go           # Root command
│   └── run.go            # Run command
├── pkg/
│   ├── types/            # Core type definitions
│   │   └── types.go      # Workflow, Task, TaskResult types
│   ├── parser/           # YAML parser
│   │   └── parser.go     # Workflow file parsing
│   ├── executor/         # Task executors
│   │   └── executor.go   # Shell, HTTP, file executors
│   └── workflow/         # Workflow engine
│       └── engine.go     # Orchestration, dependencies, parallelism
├── examples/             # Example workflows
└── main.go              # Entry point
```

## Features in Detail

### Task Dependencies

TaskFlow automatically resolves task dependencies and executes them in the correct order. Circular dependencies are detected and reported.

### Parallel Execution

Tasks marked with `parallel: true` that have the same dependencies will run concurrently, significantly reducing execution time.

### Retry Mechanism

Failed tasks can be automatically retried with configurable delays between attempts. Use `continue_on_error: true` to prevent workflow termination on failure.

### Variable Substitution

Variables can be defined at the workflow level and referenced in any task field using `${variable_name}` syntax.

### Conditional Logic

Tasks can be conditionally executed based on variable values. Supports `==` and `!=` operators.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GPL-3.0 License - see the LICENSE file for details.

## Author

Max Base

## Links

- GitHub: https://github.com/BaseMax/go-taskflow
- Issues: https://github.com/BaseMax/go-taskflow/issues

