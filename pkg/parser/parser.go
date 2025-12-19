package parser

import (
	"os"

	"github.com/BaseMax/go-taskflow/pkg/types"
	"gopkg.in/yaml.v3"
)

// ParseWorkflowFile parses a YAML workflow file
func ParseWorkflowFile(filePath string) (*types.Workflow, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var wf types.Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}

	return &wf, nil
}

// ParseWorkflowString parses a YAML workflow from a string
func ParseWorkflowString(yamlContent string) (*types.Workflow, error) {
	var wf types.Workflow
	if err := yaml.Unmarshal([]byte(yamlContent), &wf); err != nil {
		return nil, err
	}

	return &wf, nil
}
