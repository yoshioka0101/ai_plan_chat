package validation

import (
	"fmt"

	"github.com/google/uuid"
)

func ValidationTaskID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid task ID format: %w", err)
	}
	return nil
}
