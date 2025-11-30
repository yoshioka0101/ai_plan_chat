package validation

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ValidationTaskID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid task ID format: %w", err)
	}
	return nil
}

// ValidateTaskTitle はタスクタイトルの検証を行います
func ValidateTaskTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(title) > 500 {
		return fmt.Errorf("title must be 500 characters or less")
	}
	return nil
}

// ValidateTaskStatus はタスクステータスの検証を行います
func ValidateTaskStatus(status string) error {
	validStatuses := []string{"todo", "in_progress", "done", ""}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s, must be one of: todo, in_progress, done", status)
}

// ValidateTaskDueAt は期限日の検証を行います
func ValidateTaskDueAt(dueAt *time.Time) error {
	if dueAt == nil {
		return nil
	}
	now := time.Now()
	if dueAt.Before(now) {
		return fmt.Errorf("due date cannot be in the past")
	}
	return nil
}

// ValidateCreateTaskRequest はタスク作成リクエストの検証を行います
func ValidateCreateTaskRequest(title string, description *string, dueAt *time.Time, status string) error {
	if err := ValidateTaskTitle(title); err != nil {
		return err
	}
	if err := ValidateTaskStatus(status); err != nil {
		return err
	}
	if err := ValidateTaskDueAt(dueAt); err != nil {
		return err
	}
	return nil
}

// ValidateUpdateTaskRequest はタスク更新リクエストの検証を行います
func ValidateUpdateTaskRequest(title string, description *string, dueAt *time.Time, status string) error {
	return ValidateCreateTaskRequest(title, description, dueAt, status)
}

// ValidateEditTaskRequest はタスク部分更新リクエストの検証を行います
func ValidateEditTaskRequest(title *string, description *string, dueAt *time.Time, status *string) error {
	if title != nil {
		if err := ValidateTaskTitle(*title); err != nil {
			return err
		}
	}
	if status != nil {
		if err := ValidateTaskStatus(*status); err != nil {
			return err
		}
	}
	if err := ValidateTaskDueAt(dueAt); err != nil {
		return err
	}
	return nil
}
