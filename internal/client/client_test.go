package client

import (
	"context"
	"testing"
	"time"

	"github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

func TestClient_ProcessTask(t *testing.T) {
	c, err := NewClient(
		WithAPIKey("test-api-key"),
		WithRateLimit(100*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	task := &models.Task{
		Name:   "test-task",
		Prompt: "test prompt",
		Images: []string{"test.jpg"},
		Model:  "gpt-4o-image-vip",
	}

	ctx := context.Background()
	result, err := c.ProcessTask(ctx, task)
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success to be true, got false")
	}

	if result.TaskName != task.Name {
		t.Errorf("Expected task name %s, got %s", task.Name, result.TaskName)
	}
}

func TestClient_ProcessTasksAsync(t *testing.T) {
	c, err := NewClient(
		WithAPIKey("test-api-key"),
		WithRateLimit(100*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tasks := []*models.Task{
		{
			Name:   "test-task-1",
			Prompt: "test prompt 1",
			Images: []string{"test1.jpg"},
			Model:  "gpt-4o-image-vip",
		},
		{
			Name:   "test-task-2",
			Prompt: "test prompt 2",
			Images: []string{"test2.jpg"},
			Model:  "gpt-4o-image-vip",
		},
	}

	ctx := context.Background()
	resultChan := c.ProcessTasksAsync(ctx, tasks)

	count := 0
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Task failed: %v", result.Error)
		}
		count++
	}

	if count != len(tasks) {
		t.Errorf("Expected %d results, got %d", len(tasks), count)
	}
}
