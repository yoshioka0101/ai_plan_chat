package entity

import (
	"encoding/json"
	"time"
)

// ResourceType はリソースタイプ
type ResourceType string

const (
	ResourceTypeTask ResourceType = "task"
	// 今後の実装で対応する
	//ResourceTypeEvent  ResourceType = "event"
	//ResourceTypeWallet ResourceType = "wallet"
)

// ItemStatus はアイテムのステータス
type ItemStatus string

const (
	ItemStatusPending ItemStatus = "pending"
	ItemStatusCreated ItemStatus = "created"
)

// InterpretationItem はAI解釈アイテム（レビュー対象）
type InterpretationItem struct {
	ID               string
	InterpretationID string
	ItemIndex        int
	ResourceType     ResourceType
	ResourceID       *string
	Status           ItemStatus
	Data             json.RawMessage
	OriginalData     json.RawMessage
	ReviewedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TaskData はタスクアイテムのデータ構造
type TaskData struct {
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}
