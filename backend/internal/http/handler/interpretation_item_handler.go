package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	apperrors "github.com/yoshioka0101/ai_plan_chat/internal/http/errors"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/presenter"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

// InterpretationItemHandler はinterpretation itemsエンドポイントのハンドラー
type InterpretationItemHandler struct {
	itemUseCase interfaces.InterpretationItemUseCase
	presenter   *presenter.InterpretationItemPresenter
}

// NewInterpretationItemHandler はInterpretationItemHandlerを作成します
func NewInterpretationItemHandler(itemUseCase interfaces.InterpretationItemUseCase) *InterpretationItemHandler {
	return &InterpretationItemHandler{
		itemUseCase: itemUseCase,
		presenter:   presenter.NewInterpretationItemPresenter(),
	}
}

// GetInterpretationItemsByInterpretationID は指定したAI解釈IDに紐づくアイテム一覧を取得します (GET /interpretations/:id/items)
func (h *InterpretationItemHandler) GetInterpretationItemsByInterpretationID(c *gin.Context) {
	interpretationID := c.Param("id")

	// アイテム一覧を取得
	items, err := h.itemUseCase.GetItems(c.Request.Context(), interpretationID)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to get items: "+err.Error())
		return
	}

	// APIレスポンス型に変換
	apiItems, err := h.presenter.ConvertToAPIItems(items)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Failed to convert items: "+err.Error())
		return
	}

	response := api.InterpretationItemsResponse{
		Items: apiItems,
	}

	c.JSON(http.StatusOK, response)
}

// GetInterpretationItem はIDでアイテムを取得します (GET /interpretation-items/:id)
func (h *InterpretationItemHandler) GetInterpretationItem(c *gin.Context) {
	itemID := c.Param("id")

	// アイテムを取得
	item, err := h.itemUseCase.GetItem(c.Request.Context(), itemID)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrNotFound, "Item with id "+itemID+" not found")
		return
	}

	// APIレスポンス型に変換
	apiItem, err := h.presenter.ConvertToAPIItem(item)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Failed to convert item: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, apiItem)
}

// UpdateInterpretationItem はアイテムのdata内容を編集します (PATCH /interpretation-items/:id)
func (h *InterpretationItemHandler) UpdateInterpretationItem(c *gin.Context) {
	itemID := c.Param("id")

	var req api.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInvalidRequest, err.Error())
		return
	}

	// dataをJSONバイト列に変換
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInvalidRequest, "Invalid data format: "+err.Error())
		return
	}

	// アイテムを更新
	item, err := h.itemUseCase.UpdateItem(c.Request.Context(), itemID, dataBytes)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to update item: "+err.Error())
		return
	}

	// APIレスポンス型に変換
	apiItem, err := h.presenter.ConvertToAPIItem(item)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInternalServer, "Failed to convert item: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, apiItem)
}

// ApproveInterpretationItem はアイテムを承認してリソース（タスク等）を作成します (POST /interpretation-items/:id/approve)
func (h *InterpretationItemHandler) ApproveInterpretationItem(c *gin.Context) {
	itemID := c.Param("id")

	// 認証済みユーザーをcontextに埋め込む（UseCaseが参照するため）
	userID, ok := c.Get("user_id")
	if !ok {
		apperrors.RespondWithError(c, apperrors.ErrUnauthorized, "User not authenticated")
		return
	}
	ctx := c.Request.Context()
	ctx = contextWithUserID(ctx, userID.(string))

	// アイテムを承認
	resourceID, err := h.itemUseCase.ApproveItem(ctx, itemID)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to approve item: "+err.Error())
		return
	}

	// レスポンスを作成
	resourceUUID, _ := uuid.Parse(resourceID)
	response := api.ApproveItemResponse{
		ResourceId: resourceUUID,
	}

	c.JSON(http.StatusOK, response)
}

// ApproveMultipleItems は複数のアイテムを一括承認します (POST /interpretations/:id/approve-items)
func (h *InterpretationItemHandler) ApproveMultipleItems(c *gin.Context) {
	var req api.ApproveMultipleItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondWithError(c, apperrors.ErrInvalidRequest, err.Error())
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		apperrors.RespondWithError(c, apperrors.ErrUnauthorized, "User not authenticated")
		return
	}
	ctx := c.Request.Context()
	ctx = contextWithUserID(ctx, userID.(string))

	// UUIDを文字列に変換
	itemIDs := make([]string, len(req.ItemIds))
	for i, itemID := range req.ItemIds {
		itemIDs[i] = itemID.String()
	}

	// 複数アイテムを承認
	resourceIDs, err := h.itemUseCase.ApproveMultipleItems(ctx, itemIDs)
	if err != nil {
		apperrors.RespondWithError(c, apperrors.ErrDatabaseError, "Failed to approve items: "+err.Error())
		return
	}

	// レスポンスを作成（string to UUID map）
	apiResourceIDs := make(map[string]uuid.UUID)
	for itemID, resourceID := range resourceIDs {
		resourceUUID, _ := uuid.Parse(resourceID)
		apiResourceIDs[itemID] = resourceUUID
	}

	response := api.ApproveMultipleItemsResponse{
		ResourceIds: apiResourceIDs,
	}

	c.JSON(http.StatusOK, response)
}

// contextWithUserID はusecaseで参照するためのユーザーIDをContextに設定します
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}
