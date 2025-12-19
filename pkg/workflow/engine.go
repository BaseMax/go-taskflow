package workflow

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BaseMax/go-taskflow/pkg/executor"
	"github.com/BaseMax/go-taskflow/pkg/types"
)

// Engine orchestrates workflow execution
type Engine struct {
	workflow *types.Workflow
	executor *executor.Executor
	results  map[string]*types.TaskResult
	mu       sync.Mutex
}

// NewEngine creates a new workflow engine
func NewEngine(wf *types.Workflow) *Engine {
	return &Engine{
		workflow: wf,
		executor: executor.NewExecutor(wf.Variables),
		results:  make(map[string]*types.TaskResult),
	}
}

// Run executes the workflow
func (e *Engine) Run(ctx context.Context) ([]*types.TaskResult, error) {
	// Build dependency graph
	taskMap := make(map[string]*types.Task)
	for i := range e.workflow.Tasks {
		taskMap[e.workflow.Tasks[i].Name] = &e.workflow.Tasks[i]
	}

	// Execute tasks respecting dependencies
	executed := make(map[string]bool)
	var allResults []*types.TaskResult

	for len(executed) < len(e.workflow.Tasks) {
		// Find tasks that can be executed (all dependencies met)
		var readyTasks []types.Task
		for _, task := range e.workflow.Tasks {
			if executed[task.Name] {
				continue
			}

			if e.dependenciesMet(task, executed) {
				readyTasks = append(readyTasks, task)
			}
		}

		if len(readyTasks) == 0 {
			// Check if there are remaining tasks with unmet dependencies
			if len(executed) < len(e.workflow.Tasks) {
				var unexecuted []string
				for _, task := range e.workflow.Tasks {
					if !executed[task.Name] {
						unexecuted = append(unexecuted, task.Name)
					}
				}
				return allResults, fmt.Errorf("circular dependency or missing task dependencies: %v", unexecuted)
			}
			break
		}

		// Group tasks for parallel execution
		var parallelTasks, sequentialTasks []types.Task
		for _, task := range readyTasks {
			if task.Parallel {
				parallelTasks = append(parallelTasks, task)
			} else {
				sequentialTasks = append(sequentialTasks, task)
			}
		}

		// Execute parallel tasks
		if len(parallelTasks) > 0 {
			results := e.executeParallel(ctx, parallelTasks)
			for _, result := range results {
				allResults = append(allResults, result)
				executed[result.TaskName] = true
				e.storeResult(result)
			}
		}

		// Execute sequential tasks
		for _, task := range sequentialTasks {
			result := e.executeTask(ctx, task)
			allResults = append(allResults, result)
			executed[result.TaskName] = true
			e.storeResult(result)

			// Stop if task failed and continue_on_error is false
			if !result.Success && !task.ContinueOnError {
				return allResults, fmt.Errorf("task %s failed: %v", task.Name, result.Error)
			}
		}
	}

	return allResults, nil
}

// executeTask executes a single task with retry logic
func (e *Engine) executeTask(ctx context.Context, task types.Task) *types.TaskResult {
	// Check condition
	if task.Condition != "" && !e.evaluateCondition(task.Condition) {
		return &types.TaskResult{
			TaskName:  task.Name,
			Success:   true,
			Output:    "Skipped due to condition",
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}
	}

	// Set timeout if specified
	taskCtx := ctx
	if task.Timeout > 0 {
		var cancel context.CancelFunc
		taskCtx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	// Execute with retries
	maxAttempts := 1
	delay := time.Duration(0)
	if task.Retry.MaxAttempts > 0 {
		maxAttempts = task.Retry.MaxAttempts
		delay = task.Retry.Delay
	}

	var result *types.TaskResult
	var err error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 && delay > 0 {
			time.Sleep(delay)
		}

		result, err = e.executor.Execute(taskCtx, task)
		if err == nil {
			break
		}

		if attempt < maxAttempts-1 {
			fmt.Printf("Task %s failed (attempt %d/%d), retrying...\n", task.Name, attempt+1, maxAttempts)
		}
	}

	return result
}

// executeParallel executes multiple tasks in parallel
func (e *Engine) executeParallel(ctx context.Context, tasks []types.Task) []*types.TaskResult {
	var wg sync.WaitGroup
	results := make([]*types.TaskResult, len(tasks))

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t types.Task) {
			defer wg.Done()
			results[index] = e.executeTask(ctx, t)
		}(i, task)
	}

	wg.Wait()
	return results
}

// dependenciesMet checks if all dependencies for a task are satisfied
func (e *Engine) dependenciesMet(task types.Task, executed map[string]bool) bool {
	for _, dep := range task.DependsOn {
		if !executed[dep] {
			return false
		}
		// Also check if the dependency was successful
		e.mu.Lock()
		result, exists := e.results[dep]
		e.mu.Unlock()
		if exists && !result.Success {
			return false
		}
	}
	return true
}

// storeResult stores a task result
func (e *Engine) storeResult(result *types.TaskResult) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.results[result.TaskName] = result
}

// evaluateCondition evaluates a simple condition
// Supports basic expressions like: ${VAR} == "value", ${VAR} != "value"
func (e *Engine) evaluateCondition(condition string) bool {
	// Replace variables
	condition = e.replaceVariables(condition)

	// Simple evaluation - check for common patterns
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(strings.Trim(parts[1], "\"'"))
			return left == right
		}
	}

	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(strings.Trim(parts[1], "\"'"))
			return left != right
		}
	}

	// If condition is just a variable, check if it's truthy
	return condition != "" && condition != "false" && condition != "0"
}

// replaceVariables replaces ${VAR} patterns with actual variable values
func (e *Engine) replaceVariables(input string) string {
	result := input
	for key, value := range e.workflow.Variables {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
	}
	// Also check task results
	e.mu.Lock()
	for taskName, taskResult := range e.results {
		if taskResult.Success {
			result = strings.ReplaceAll(result, "${"+taskName+".output}", taskResult.Output)
		}
	}
	e.mu.Unlock()
	return result
}
