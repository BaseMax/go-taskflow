# TaskFlow Quick Start Guide

This guide will help you get started with TaskFlow in just a few minutes.

## Installation

### Option 1: Build from Source

```bash
git clone https://github.com/BaseMax/go-taskflow
cd go-taskflow
make build
```

The `taskflow` binary will be created in the current directory.

### Option 2: Install to System

```bash
make install
```

This will install TaskFlow to `/usr/local/bin/` (requires sudo).

## Your First Workflow

Create a file named `hello.yaml`:

```yaml
name: Hello World
description: My first TaskFlow workflow

tasks:
  - name: greet
    type: shell
    command: echo "Hello, TaskFlow!"
```

Run it:

```bash
./taskflow run hello.yaml
```

## Common Workflows

### Sequential Tasks with Dependencies

```yaml
name: Build and Deploy
tasks:
  - name: install
    type: shell
    command: npm install

  - name: build
    type: shell
    command: npm run build
    depends_on:
      - install

  - name: deploy
    type: shell
    command: npm run deploy
    depends_on:
      - build
```

### Parallel Execution

```yaml
name: Run Tests in Parallel
tasks:
  - name: setup
    type: shell
    command: echo "Setting up..."

  - name: unit-tests
    type: shell
    command: npm run test:unit
    depends_on: [setup]
    parallel: true

  - name: integration-tests
    type: shell
    command: npm run test:integration
    depends_on: [setup]
    parallel: true

  - name: e2e-tests
    type: shell
    command: npm run test:e2e
    depends_on: [setup]
    parallel: true
```

### Using Variables

```yaml
name: Deploy with Variables
variables:
  app_name: "myapp"
  environment: "production"
  version: "1.0.0"

tasks:
  - name: deploy
    type: shell
    command: |
      echo "Deploying ${app_name} v${version} to ${environment}"
      kubectl apply -f deployment.yaml
```

### HTTP Requests

```yaml
name: API Testing
tasks:
  - name: health-check
    type: http
    url: https://api.example.com/health
    method: GET

  - name: create-user
    type: http
    url: https://api.example.com/users
    method: POST
    headers:
      Content-Type: application/json
    body: '{"name": "John Doe", "email": "john@example.com"}'
```

### File Operations

```yaml
name: Backup Configuration
variables:
  config_file: "/etc/app/config.json"
  backup_dir: "/backups"

tasks:
  - name: create-backup-dir
    type: shell
    command: mkdir -p ${backup_dir}

  - name: backup-config
    type: file
    file_action: copy
    source_path: ${config_file}
    dest_path: ${backup_dir}/config.backup.json
    depends_on:
      - create-backup-dir
```

### Retry Failed Tasks

```yaml
name: Flaky Operation
tasks:
  - name: api-call
    type: http
    url: https://api.example.com/data
    method: GET
    retry:
      max_attempts: 3
      delay: 5s
```

### Conditional Execution

```yaml
name: Environment-Specific Tasks
variables:
  env: "staging"

tasks:
  - name: deploy-staging
    type: shell
    command: ./deploy-staging.sh
    condition: ${env} == "staging"

  - name: deploy-production
    type: shell
    command: ./deploy-production.sh
    condition: ${env} == "production"
```

### Task Timeout

```yaml
name: Long Running Task
tasks:
  - name: backup
    type: shell
    command: ./backup-database.sh
    timeout: 30m
```

### Continue on Error

```yaml
name: Cleanup Pipeline
tasks:
  - name: risky-operation
    type: shell
    command: ./might-fail.sh
    continue_on_error: true

  - name: cleanup
    type: shell
    command: ./cleanup.sh
    depends_on:
      - risky-operation
```

## Real-World Example: CI/CD Pipeline

See `examples/ci-pipeline.yaml` for a complete CI/CD pipeline example.

```bash
./taskflow run examples/ci-pipeline.yaml
```

## Getting Help

- View all available commands: `./taskflow --help`
- View command help: `./taskflow run --help`
- Check version: `./taskflow --version`

## Next Steps

- Explore more examples in the `examples/` directory
- Read the full documentation in [README.md](README.md)
- Check out [CONTRIBUTING.md](CONTRIBUTING.md) to contribute

## Tips

1. **Start Simple**: Begin with basic shell commands and gradually add complexity
2. **Use Dependencies**: Structure workflows with proper task dependencies
3. **Leverage Parallelism**: Speed up workflows by running independent tasks in parallel
4. **Handle Failures**: Use retry logic and `continue_on_error` for robust workflows
5. **Test Locally**: Test workflows locally before deploying to production

Happy automating! 🚀
