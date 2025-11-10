package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	apperrors "github.com/yoshioka0101/ai_plan_chat/internal/http/errors"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
	"github.com/yoshioka0101/ai_plan_chat/internal/service"
)

// InterpretationHandler はAI解釈エンドポイントのハンドラー
type InterpretationHandler struct {
	geminiService          *service.GeminiService
	interpretationRepo     interfaces.InterpretationRepository
}

// NewInterpretationHandler はInterpretationHandlerを作成します
func NewInterpretationHandler(geminiService *service.GeminiService, interpretationRepo interfaces.InterpretationRepository) *InterpretationHandler {
	return &InterpretationHandler{
		geminiService:      geminiService,
		interpretationRepo: interpretationRepo,
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
	result, err := h.geminiService.InterpretInput(c.Request.Context(), inputText)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrAIInterpretationError, "Failed to interpret input: "+err.Error())
		return
	}

	// TODO: 実際のユーザーIDは認証トークンから取得する
	// 現時点では仮のユーザーIDを使用
	userID := uuid.New().String()
	interpretationID := uuid.New().String()

	// Entity型でデータベースに保存
	entityInterpretation := &entity.AIInterpretation{
		ID:                 interpretationID,
		UserID:             userID,
		InputText:          inputText,
		Result:             *result,
		AIModel:            h.geminiService.ModelName(),
		AIPromptTokens:     ptrInt(len(inputText) / 4), // 概算
		AICompletionTokens: ptrInt(100),                // 概算
	}

	// データベースに保存
	if err := h.interpretationRepo.CreateInterpretation(c.Request.Context(), entityInterpretation); err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to save interpretation: "+err.Error())
		return
	}

	// レスポンスを作成
	interpretationIDUUID, _ := uuid.Parse(interpretationID)
	userIDUUID, _ := uuid.Parse(userID)

	interpretation := buildAIInterpretation(
		interpretationIDUUID,
		userIDUUID,
		inputText,
		result,
		h.geminiService.ModelName(),
		entityInterpretation.CreatedAt,
	)

	interpretationType := convertToResponseType(result.Type)
	response := api.InterpretationResponse{
		Type:           interpretationType,
		Interpretation: interpretation,
		Message:        nil,
	}

	c.JSON(http.StatusOK, response)
}

// ListInterpretations はAI解釈履歴を取得します
func (h *InterpretationHandler) ListInterpretations(c *gin.Context) {
	// TODO: 実際のユーザーIDは認証トークンから取得する
	// 現時点では仮のユーザーIDを使用（全件取得のため適当なID）
	userID := "00000000-0000-0000-0000-000000000000"

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
			Deadline *time.Time                                              `json:"deadline,omitempty"`
			Priority *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
			Tags     *[]string                                               `json:"tags,omitempty"`
		} `json:"metadata,omitempty"`
		Title *string                                     `json:"title,omitempty"`
		Type  *api.AIInterpretationStructuredResultType `json:"type,omitempty"`
	}{
		Title:       ptrString(result.Title),
		Description: ptrStringIfNotEmpty(result.Description),
	}

	// Metadataの処理（Todoのみ）
	metadata := &struct {
		Deadline *time.Time                                              `json:"deadline,omitempty"`
		Priority *api.AIInterpretationStructuredResultMetadataPriority `json:"priority,omitempty"`
		Tags     *[]string                                               `json:"tags,omitempty"`
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
