package entity

import (
	"encoding/json"
	"time"
)

// ItemType represents the type of item
type ItemType string

const (
	ItemTypeTodo     ItemType = "todo"
	ItemTypeEvent    ItemType = "event"
	ItemTypeReminder ItemType = "reminder"
	ItemTypeNote     ItemType = "note"
	ItemTypeExpense  ItemType = "expense"
)

// Metadata contains AI-extracted information from user input
type Metadata struct {
	// Common fields
	Deadline  *time.Time `json:"deadline,omitempty"`
	Priority  *string    `json:"priority,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Location  *string    `json:"location,omitempty"`

	// Expense-specific fields
	Amount   *float64 `json:"amount,omitempty"`
	Currency *string  `json:"currency,omitempty"`
	Category *string  `json:"category,omitempty"`

	// Event-specific fields
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`

	// Additional custom fields
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// Item represents a user input interpreted by AI
type Item struct {
	ID                string
	UserID            string
	Title             string
	Type              ItemType
	Metadata          Metadata
	OriginalInput     string
	AIModel           string
	AIPromptTokens    *int
	AICompletionTokens *int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// MarshalMetadata converts Metadata to JSON bytes
func (m *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(*m)
}

// UnmarshalMetadata converts JSON bytes to Metadata
func UnmarshalMetadata(data []byte) (*Metadata, error) {
	var m Metadata
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Validate validates the Item
func (i *Item) Validate() error {
	if i.UserID == "" {
		return ErrEmptyUserID
	}
	if i.Title == "" {
		return ErrEmptyTitle
	}
	if i.OriginalInput == "" {
		return ErrEmptyOriginalInput
	}
	if i.AIModel == "" {
		return ErrEmptyAIModel
	}
	if !isValidItemType(i.Type) {
		return ErrInvalidItemType
	}
	return nil
}

// isValidItemType checks if the given item type is valid
func isValidItemType(t ItemType) bool {
	switch t {
	case ItemTypeTodo, ItemTypeEvent, ItemTypeReminder, ItemTypeNote, ItemTypeExpense:
		return true
	default:
		return false
	}
}
