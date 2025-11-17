package entity

import "time"

// InterpretationType はAI解釈のタイプ
type InterpretationType string

const (
	InterpretationTypeTodo InterpretationType = "todo"
)

// InterpretationResult はAIによる解釈結果
type InterpretationResult struct {
	Type        InterpretationType
	Title       string
	Description string
	Metadata    InterpretationMetadata
}

// InterpretationMetadata は解釈結果のメタデータ（Todoのみ）
type InterpretationMetadata struct {
	// Todo関連フィールド
	Deadline *time.Time `json:"deadline,omitempty"`
	Priority *string    `json:"priority,omitempty"`
	Tags     []string   `json:"tags,omitempty"`

	// 追加のカスタムフィールド
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// AIInterpretation はAI解釈の完全な情報（DB保存用）
type AIInterpretation struct {
	ID                 string
	UserID             string
	InputText          string
	Result             InterpretationResult
	AIModel            string
	AIPromptTokens     *int
	AICompletionTokens *int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
