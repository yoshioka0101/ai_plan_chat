package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
)

// InterpretationHandler はAI解釈エンドポイントのハンドラー
type InterpretationHandler struct{}

// NewInterpretationHandler はInterpretationHandlerを作成します
func NewInterpretationHandler() *InterpretationHandler {
	return &InterpretationHandler{}
}

// CreateInterpretation は自然言語入力からAI解釈を作成します（デモ実装）
func (h *InterpretationHandler) CreateInterpretation(c *gin.Context) {
	var req api.CreateInterpretationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Code:    ptrString("invalid_request"),
			Message: ptrString(err.Error()),
		})
		return
	}

	// デモ実装: 入力テキストをそのままエコーして、簡単な解析結果を返す
	inputText := req.InputText

	// 簡単なキーワードベースの判定
	interpretationType := determineType(inputText)

	// デモレスポンスを作成
	now := time.Now()
	uuidVal, _ := uuid.Parse(uuid.New().String())
	demoUserUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")

	tags := []string{"demo", "test"}
	interpretation := api.AIInterpretation{
		Id:        openapi_types.UUID(uuidVal),
		UserId:    openapi_types.UUID(demoUserUUID),
		InputText: inputText,
		StructuredResult: struct {
			Description *string `json:"description,omitempty"`
			Metadata    *struct {
				Amount    *float32   `json:"amount,omitempty"`
				Category  *string    `json:"category,omitempty"`
				Currency  *string    `json:"currency,omitempty"`
				Deadline  *time.Time `json:"deadline,omitempty"`
				EndTime   *time.Time `json:"end_time,omitempty"`
				Location  *string    `json:"location,omitempty"`
				Priority  *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
				StartTime *time.Time `json:"start_time,omitempty"`
				Tags      *[]string  `json:"tags,omitempty"`
			} `json:"metadata,omitempty"`
			Title *string                                     `json:"title,omitempty"`
			Type  *api.AIInterpretationStructuredResultType `json:"type,omitempty"`
		}{
			Title:       ptrString(inputText),
			Description: ptrString("これはデモ実装です。入力: " + inputText),
			Metadata: &struct {
				Amount    *float32   `json:"amount,omitempty"`
				Category  *string    `json:"category,omitempty"`
				Currency  *string    `json:"currency,omitempty"`
				Deadline  *time.Time `json:"deadline,omitempty"`
				EndTime   *time.Time `json:"end_time,omitempty"`
				Location  *string    `json:"location,omitempty"`
				Priority  *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
				StartTime *time.Time `json:"start_time,omitempty"`
				Tags      *[]string  `json:"tags,omitempty"`
			}{
				Tags: &tags,
			},
		},
		AiModel:            "demo-model-v1",
		AiPromptTokens:     ptrInt(len(inputText)),
		AiCompletionTokens: ptrInt(50),
		CreatedAt:          now,
	}

	response := api.InterpretationResponse{
		Type:           interpretationType,
		Interpretation: interpretation,
		Message:        nil,
	}

	c.JSON(http.StatusOK, response)
}

// ListInterpretations はAI解釈履歴を取得します（デモ実装）
func (h *InterpretationHandler) ListInterpretations(c *gin.Context) {
	// デモ実装: 空の配列を返す
	c.JSON(http.StatusOK, gin.H{
		"interpretations": []api.AIInterpretation{},
		"total":           0,
		"limit":           20,
		"offset":          0,
	})
}

// GetInterpretation は特定のAI解釈を取得します（デモ実装）
func (h *InterpretationHandler) GetInterpretation(c *gin.Context) {
	id := c.Param("id")

	// デモ実装: 404を返す
	c.JSON(http.StatusNotFound, api.ErrorResponse{
		Code:    ptrString("not_found"),
		Message: ptrString("Interpretation with id " + id + " not found"),
	})
}

// determineType は入力テキストから簡単にタイプを判定します（デモ実装）
func determineType(text string) api.InterpretationResponseType {
	// 簡単なキーワードマッチング
	keywords := map[string]api.InterpretationResponseType{
		"買う":     api.InterpretationResponseTypeTodo,
		"購入":     api.InterpretationResponseTypeTodo,
		"やる":     api.InterpretationResponseTypeTodo,
		"タスク":    api.InterpretationResponseTypeTodo,
		"会議":     api.InterpretationResponseTypeEvent,
		"ミーティング": api.InterpretationResponseTypeEvent,
		"予定":     api.InterpretationResponseTypeEvent,
		"イベント":   api.InterpretationResponseTypeEvent,
		"思い出":    api.InterpretationResponseTypeReminder,
		"リマインド":  api.InterpretationResponseTypeReminder,
		"忘れない":   api.InterpretationResponseTypeReminder,
		"メモ":     api.InterpretationResponseTypeNote,
		"覚えて":    api.InterpretationResponseTypeNote,
		"記録":     api.InterpretationResponseTypeNote,
		"円":      api.InterpretationResponseTypeExpense,
		"ドル":     api.InterpretationResponseTypeExpense,
		"支払":     api.InterpretationResponseTypeExpense,
		"経費":     api.InterpretationResponseTypeExpense,
	}

	for keyword, typeVal := range keywords {
		if containsKeyword(text, keyword) {
			return typeVal
		}
	}

	return api.InterpretationResponseTypeUnknown
}

// containsKeyword はテキストにキーワードが含まれているかチェックします
func containsKeyword(text, keyword string) bool {
	// 簡易的な部分文字列マッチング
	for i := 0; i <= len(text)-len(keyword); i++ {
		if text[i:i+len(keyword)] == keyword {
			return true
		}
	}
	return false
}

// ptrString はstringのポインタを返すヘルパー関数
func ptrString(s string) *string {
	return &s
}

// ptrInt はintのポインタを返すヘルパー関数
func ptrInt(i int) *int {
	return &i
}
