package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	apperrors "github.com/yoshioka0101/ai_plan_chat/internal/http/errors"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
	"github.com/yoshioka0101/ai_plan_chat/internal/service"
)

// InterpretationHandler はAI解釈エンドポイントのハンドラー
type InterpretationHandler struct {
	geminiService          *service.GeminiService
	interpretationRepo     interfaces.InterpretationRepository
	interpretationItemRepo interfaces.InterpretationItemRepository
}

// NewInterpretationHandler はInterpretationHandlerを作成します
func NewInterpretationHandler(geminiService *service.GeminiService, interpretationRepo interfaces.InterpretationRepository, interpretationItemRepo interfaces.InterpretationItemRepository) *InterpretationHandler {
	return &InterpretationHandler{
		geminiService:          geminiService,
		interpretationRepo:     interpretationRepo,
		interpretationItemRepo: interpretationItemRepo,
	}
}

// CreateInterpretation は自然言語入力からAI解釈を作成します
func (h *InterpretationHandler) CreateInterpretation(c *gin.Context) {
	var req api.CreateInterpretationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInvalidRequest, err.Error())
		return
	}

	inputText := req.InputText

	// Geminiサービスの存在チェック
	if h.geminiService == nil {
		apperrors.RespondWithError(c, apperrors.ErrConfigurationError, "AI service is not configured")
		return
	}

	// Gemini APIで解析
	aiResult, err := h.geminiService.InterpretInput(c.Request.Context(), inputText)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrAIInterpretationError, "Failed to interpret input: "+err.Error())
		return
	}

	// 認証ミドルウェアから設定されたユーザーIDを取得
	userIDValue, exists := c.Get("user_id")
	if !exists {
		apperrors.RespondWithError(c, apperrors.ErrUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDValue.(string)
	if !ok {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Invalid user ID format")
		return
	}

	interpretationID := uuid.New().String()

	// Entity型でデータベースに保存
	entityInterpretation := &entity.AIInterpretation{
		ID:                 interpretationID,
		UserID:             userID,
		InputText:          inputText,
		Result:             *aiResult.Result,
		OriginalResult:     aiResult.OriginalJSON,
		AIModel:            h.geminiService.ModelName(),
		AIPromptTokens:     ptrInt(len(inputText) / 4), // 概算
		AICompletionTokens: ptrInt(100),                // 概算
	}

	// データベースに保存
	if err := h.interpretationRepo.CreateInterpretation(c.Request.Context(), entityInterpretation); err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to save interpretation: "+err.Error())
		return
	}

	if h.interpretationItemRepo == nil {
		apperrors.RespondWithError(c, apperrors.ErrConfigurationError, "Item repository is not configured")
		return
	}

	items, err := buildInterpretationItems(interpretationID, aiResult.Result, aiResult.OriginalJSON)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Failed to prepare interpretation items: "+err.Error())
		return
	}

	if len(items) > 0 {
		if err := h.interpretationItemRepo.CreateItems(c.Request.Context(), items); err != nil {
			apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to save interpretation items: "+err.Error())
			return
		}
	}

	// レスポンスを作成
	interpretationIDUUID, _ := uuid.Parse(interpretationID)
	userIDUUID, _ := uuid.Parse(userID)

	interpretation := buildAIInterpretation(
		interpretationIDUUID,
		userIDUUID,
		inputText,
		aiResult.Result,
		h.geminiService.ModelName(),
		entityInterpretation.CreatedAt,
	)

	interpretationType := convertToResponseType(aiResult.Result.Type)
	response := api.InterpretationResponse{
		Type:           interpretationType,
		Interpretation: interpretation,
		Message:        nil,
	}

	c.JSON(http.StatusOK, response)
}

// ListInterpretations はAI解釈履歴を取得します
func (h *InterpretationHandler) ListInterpretations(c *gin.Context) {
	// 認証ミドルウェアから設定されたユーザーIDを取得
	userIDValue, exists := c.Get("user_id")
	if !exists {
		apperrors.RespondWithError(c, apperrors.ErrUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDValue.(string)
	if !ok {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Invalid user ID format")
		return
	}

	limit := 20
	offset := 0

	// データベースから取得
	interpretations, err := h.interpretationRepo.GetInterpretationsByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to get interpretations: "+err.Error())
		return
	}

	// APIレスポンス型に変換
	apiInterpretations := make([]api.AIInterpretation, 0, len(interpretations))
	for _, interp := range interpretations {
		interpretationIDUUID, _ := uuid.Parse(interp.ID)
		userIDUUID, _ := uuid.Parse(interp.UserID)

		apiInterp := buildAIInterpretation(
			interpretationIDUUID,
			userIDUUID,
			interp.InputText,
			&interp.Result,
			interp.AIModel,
			interp.CreatedAt,
		)
		apiInterpretations = append(apiInterpretations, apiInterp)
	}

	c.JSON(http.StatusOK, gin.H{
		"interpretations": apiInterpretations,
		"total":           len(apiInterpretations),
		"limit":           limit,
		"offset":          offset,
	})
}

// GetInterpretation は特定のAI解釈を取得します
func (h *InterpretationHandler) GetInterpretation(c *gin.Context) {
	id := c.Param("id")

	// データベースから取得
	interpretation, err := h.interpretationRepo.GetInterpretationByID(c.Request.Context(), id)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrNotFound, "Interpretation with id "+id+" not found")
		return
	}

	// APIレスポンス型に変換
	interpretationIDUUID, _ := uuid.Parse(interpretation.ID)
	userIDUUID, _ := uuid.Parse(interpretation.UserID)

	apiInterp := buildAIInterpretation(
		interpretationIDUUID,
		userIDUUID,
		interpretation.InputText,
		&interpretation.Result,
		interpretation.AIModel,
		interpretation.CreatedAt,
	)

	c.JSON(http.StatusOK, apiInterp)
}

// buildInterpretationItems はAI解釈結果からレビュー用アイテムを組み立てます
func buildInterpretationItems(interpretationID string, result *entity.InterpretationResult, originalJSON []byte) ([]*entity.InterpretationItem, error) {
	if result == nil {
		return nil, fmt.Errorf("interpretation result is nil")
	}

	taskData := entity.TaskData{
		Title: result.Title,
	}

	if desc := ptrStringIfNotEmpty(result.Description); desc != nil {
		taskData.Description = desc
	}

	if result.Metadata.Deadline != nil {
		taskData.DueAt = result.Metadata.Deadline
	}

	if result.Metadata.Priority != nil {
		taskData.Priority = result.Metadata.Priority
	}

	if len(result.Metadata.Tags) > 0 {
		taskData.Tags = result.Metadata.Tags
	}

	dataBytes, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	itemOriginal := originalJSON
	if len(itemOriginal) == 0 {
		itemOriginal = dataBytes
	}

	return []*entity.InterpretationItem{
		{
			ID:               uuid.New().String(),
			InterpretationID: interpretationID,
			ItemIndex:        0,
			ResourceType:     entity.ResourceTypeTask,
			Status:           entity.ItemStatusPending,
			Data:             dataBytes,
			OriginalData:     itemOriginal,
		},
	}, nil
}

// buildAIInterpretation はAIInterpretation構造体を構築します
func buildAIInterpretation(
	id uuid.UUID,
	userID uuid.UUID,
	inputText string,
	result *entity.InterpretationResult,
	modelName string,
	createdAt time.Time,
) api.AIInterpretation {
	structuredResult := struct {
		Description *string `json:"description,omitempty"`
		Metadata    *struct {
			Deadline *time.Time                                            `json:"deadline,omitempty"`
			Priority *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
			Tags     *[]string                                             `json:"tags,omitempty"`
		} `json:"metadata,omitempty"`
		Title *string                                   `json:"title,omitempty"`
		Type  *api.AIInterpretationStructuredResultType `json:"type,omitempty"`
	}{
		Title:       ptrString(result.Title),
		Description: ptrStringIfNotEmpty(result.Description),
	}

	// Metadataの処理（Todoのみ）
	metadata := &struct {
		Deadline *time.Time                                            `json:"deadline,omitempty"`
		Priority *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
		Tags     *[]string                                             `json:"tags,omitempty"`
	}{
		Deadline: result.Metadata.Deadline,
	}

	// タグの設定
	if len(result.Metadata.Tags) > 0 {
		metadata.Tags = &result.Metadata.Tags
	}

	// 優先度の変換
	if result.Metadata.Priority != nil {
		priority := api.AIInterpretationStructuredResultMetadataPriority(*result.Metadata.Priority)
		metadata.Priority = &priority
	}

	structuredResult.Metadata = metadata

	return api.AIInterpretation{
		Id:                 openapi_types.UUID(id),
		UserId:             openapi_types.UUID(userID),
		InputText:          inputText,
		StructuredResult:   structuredResult,
		AiModel:            modelName,
		AiPromptTokens:     ptrInt(len(inputText) / 4), // 概算
		AiCompletionTokens: ptrInt(100),                // 概算
		CreatedAt:          createdAt,
	}
}

// convertToResponseType はentityのタイプをAPIのタイプに変換します（現在はtodoのみ）
func convertToResponseType(t entity.InterpretationType) api.InterpretationResponseType {
	// 現在はtodoのみをサポート
	return api.InterpretationResponseTypeTodo
}

// ptrStringIfNotEmpty は空でない場合のみstringのポインタを返すヘルパー関数
func ptrStringIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrString はstringのポインタを返すヘルパー関数
func ptrString(s string) *string {
	return &s
}

// ptrInt はintのポインタを返すヘルパー関数
func ptrInt(i int) *int {
	return &i
}
