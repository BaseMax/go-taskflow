package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/BaseMax/go-taskflow/pkg/types"
)

const (
	// Shell command constants
	windowsShell = "cmd.exe"
	windowsFlag  = "/c"
	unixShell    = "sh"
	unixFlag     = "-c"
)

// Executor handles execution of different task types
type Executor struct {
	variables map[string]string
}

// NewExecutor creates a new task executor
func NewExecutor(variables map[string]string) *Executor {
	return &Executor{
		variables: variables,
	}
}

// Execute executes a task based on its type
func (e *Executor) Execute(ctx context.Context, task types.Task) (*types.TaskResult, error) {
	result := &types.TaskResult{
		TaskName:  task.Name,
		StartTime: time.Now(),
	}

	var err error
	var output string

	switch task.Type {
	case "shell":
		output, err = e.executeShell(ctx, task)
	case "http":
		output, err = e.executeHTTP(ctx, task)
	case "file":
		output, err = e.executeFile(ctx, task)
	default:
		err = fmt.Errorf("unknown task type: %s", task.Type)
	}

	result.EndTime = time.Now()
	result.Output = output
	result.Success = err == nil
	result.Error = err

	return result, err
}

// executeShell executes a shell command
func (e *Executor) executeShell(ctx context.Context, task types.Task) (string, error) {
	command := task.Command
	if task.Script != "" {
		command = task.Script
	}

	// Replace variables in command
	command = e.replaceVariables(command)

	// Use appropriate shell based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, windowsShell, windowsFlag, command)
	} else {
		cmd = exec.CommandContext(ctx, unixShell, unixFlag, command)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), fmt.Errorf("command failed: %w\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// executeHTTP executes an HTTP request
func (e *Executor) executeHTTP(ctx context.Context, task types.Task) (string, error) {
	url := e.replaceVariables(task.URL)
	method := task.Method
	if method == "" {
		method = "GET"
	}

	var body io.Reader
	if task.Body != "" {
		body = strings.NewReader(e.replaceVariables(task.Body))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range task.Headers {
		req.Header.Set(key, e.replaceVariables(value))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return string(respBody), fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	return string(respBody), nil
}

// executeFile executes file operations
func (e *Executor) executeFile(ctx context.Context, task types.Task) (string, error) {
	switch task.FileAction {
	case "read":
		return e.fileRead(task)
	case "write":
		return e.fileWrite(task)
	case "delete":
		return e.fileDelete(task)
	case "copy":
		return e.fileCopy(task)
	default:
		return "", fmt.Errorf("unknown file action: %s", task.FileAction)
	}
}

// fileRead reads content from a file
func (e *Executor) fileRead(task types.Task) (string, error) {
	filePath := e.replaceVariables(task.FilePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// fileWrite writes content to a file
func (e *Executor) fileWrite(task types.Task) (string, error) {
	filePath := e.replaceVariables(task.FilePath)
	content := e.replaceVariables(task.FileContent)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return fmt.Sprintf("Successfully wrote to %s", filePath), nil
}

// fileDelete deletes a file
func (e *Executor) fileDelete(task types.Task) (string, error) {
	filePath := e.replaceVariables(task.FilePath)
	err := os.Remove(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to delete file: %w", err)
	}
	return fmt.Sprintf("Successfully deleted %s", filePath), nil
}

// fileCopy copies a file
func (e *Executor) fileCopy(task types.Task) (string, error) {
	sourcePath := e.replaceVariables(task.SourcePath)
	destPath := e.replaceVariables(task.DestPath)

	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	err = os.WriteFile(destPath, sourceData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write destination file: %w", err)
	}

	return fmt.Sprintf("Successfully copied %s to %s", sourcePath, destPath), nil
}

// replaceVariables replaces ${VAR} patterns with actual variable values
func (e *Executor) replaceVariables(input string) string {
	result := input
	for key, value := range e.variables {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
	}
	return result
}
