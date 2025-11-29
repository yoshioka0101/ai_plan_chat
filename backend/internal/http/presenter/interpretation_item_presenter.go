package presenter

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
)

// InterpretationItemPresenter はInterpretationItemのレスポンス整形を担当します
type InterpretationItemPresenter struct{}

// NewInterpretationItemPresenter はInterpretationItemPresenterを作成します
func NewInterpretationItemPresenter() *InterpretationItemPresenter {
	return &InterpretationItemPresenter{}
}

// ConvertToAPIItem はentityのInterpretationItemをAPI型に変換します
func (p *InterpretationItemPresenter) ConvertToAPIItem(item *entity.InterpretationItem) (api.InterpretationItem, error) {
	// IDのパース
	id, err := uuid.Parse(item.ID)
	if err != nil {
		log.Printf("Warning: invalid UUID in database: %s, error: %v", item.ID, err)
		id = uuid.Nil
	}

	interpretationID, err := uuid.Parse(item.InterpretationID)
	if err != nil {
		log.Printf("Warning: invalid interpretation UUID in database: %s, error: %v", item.InterpretationID, err)
		interpretationID = uuid.Nil
	}

	// dataとoriginal_dataをmap[string]interface{}に変換
	var data map[string]interface{}
	if err := json.Unmarshal(item.Data, &data); err != nil {
		return api.InterpretationItem{}, err
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(item.OriginalData, &originalData); err != nil {
		return api.InterpretationItem{}, err
	}

	apiItem := api.InterpretationItem{
		Id:               types.UUID(id),
		InterpretationId: types.UUID(interpretationID),
		ItemIndex:        item.ItemIndex,
		ResourceType:     api.InterpretationItemResourceType(item.ResourceType),
		Status:           api.InterpretationItemStatus(item.Status),
		Data:             data,
		OriginalData:     originalData,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}

	// Null可能なフィールドの処理
	if item.ResourceID != nil {
		resourceID, err := uuid.Parse(*item.ResourceID)
		if err == nil {
			resourceUUID := types.UUID(resourceID)
			apiItem.ResourceId = &resourceUUID
		}
	}

	if item.ReviewedAt != nil {
		apiItem.ReviewedAt = item.ReviewedAt
	}

	return apiItem, nil
}

// ConvertToAPIItems はentityスライスをAPI型スライスに変換します
func (p *InterpretationItemPresenter) ConvertToAPIItems(items []*entity.InterpretationItem) ([]api.InterpretationItem, error) {
	apiItems := make([]api.InterpretationItem, 0, len(items))
	for _, item := range items {
		apiItem, err := p.ConvertToAPIItem(item)
		if err != nil {
			return nil, err
		}
		apiItems = append(apiItems, apiItem)
	}
	return apiItems, nil
}
