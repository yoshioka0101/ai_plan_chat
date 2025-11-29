package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/stephenafamo/bob/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

type interpretationItemRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

// NewInterpretationItemRepository は新しいInterpretationItemRepositoryを生成します
func NewInterpretationItemRepository(db bob.Executor, logger *slog.Logger) interfaces.InterpretationItemRepository {
	return &interpretationItemRepository{
		db:     db,
		logger: logger,
	}
}

// GetItemsByInterpretationID は指定したAI解釈IDに紐づくアイテム一覧を取得します
func (r *interpretationItemRepository) GetItemsByInterpretationID(ctx context.Context, interpretationID string) ([]*entity.InterpretationItem, error) {
	r.logger.InfoContext(ctx, "Repository: GetItemsByInterpretationID started",
		slog.String("interpretation_id", interpretationID),
	)

	dbItems, err := models.InterpretationItems.Query(
		sm.Where(models.InterpretationItems.Columns.InterpretationID.EQ(mysql.Arg(interpretationID))),
		sm.OrderBy(models.InterpretationItems.Columns.ItemIndex).Asc(),
	).All(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to query items",
			slog.String("interpretation_id", interpretationID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	items := make([]*entity.InterpretationItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		item, err := r.toEntity(dbItem)
		if err != nil {
			r.logger.ErrorContext(ctx, "Repository: Failed to convert item to entity",
				slog.String("item_id", dbItem.ID),
				slog.String("error", err.Error()),
			)
			continue
		}
		items = append(items, item)
	}

	r.logger.InfoContext(ctx, "Repository: GetItemsByInterpretationID completed",
		slog.String("interpretation_id", interpretationID),
		slog.Int("count", len(items)),
	)
	return items, nil
}

// GetItemByID はIDでアイテムを取得します
func (r *interpretationItemRepository) GetItemByID(ctx context.Context, id string) (*entity.InterpretationItem, error) {
	r.logger.InfoContext(ctx, "Repository: GetItemByID started",
		slog.String("item_id", id),
	)

	dbItem, err := models.InterpretationItems.Query(
		sm.Where(models.InterpretationItems.Columns.ID.EQ(mysql.Arg(id))),
	).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "Repository: Item not found",
				slog.String("item_id", id),
			)
			return nil, fmt.Errorf("item not found: %s", id)
		}
		r.logger.ErrorContext(ctx, "Repository: Failed to query item",
			slog.String("item_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	item, err := r.toEntity(dbItem)
	if err != nil {
		return nil, fmt.Errorf("failed to convert item to entity: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetItemByID completed",
		slog.String("item_id", id),
	)
	return item, nil
}

// CreateItems はアイテムをバルク作成します
func (r *interpretationItemRepository) CreateItems(ctx context.Context, items []*entity.InterpretationItem) error {
	r.logger.InfoContext(ctx, "Repository: CreateItems started",
		slog.Int("count", len(items)),
	)

	if len(items) == 0 {
		return nil
	}

	// 個別にインサート（トランザクション内で実行されることを想定）
	for _, item := range items {
		if item.ID == "" {
			item.ID = uuid.New().String()
		}
		now := time.Now()
		item.CreatedAt = now
		item.UpdatedAt = now

		setter := &models.InterpretationItemSetter{
			ID:               omit.From(item.ID),
			InterpretationID: omit.From(item.InterpretationID),
			ItemIndex:        omit.From(int32(item.ItemIndex)),
			ResourceType:     omit.From(string(item.ResourceType)),
			ResourceID:       omitnull.FromNull(null.FromPtr(item.ResourceID)),
			Status:           omit.From(string(item.Status)),
			Data:             omit.From(types.JSON[json.RawMessage]{Val: item.Data}),
			OriginalData:     omit.From(types.JSON[json.RawMessage]{Val: item.OriginalData}),
			ReviewedAt:       omitnull.FromNull(null.FromPtr(item.ReviewedAt)),
			CreatedAt:        omit.From(item.CreatedAt),
			UpdatedAt:        omit.From(item.UpdatedAt),
		}

		_, err := models.InterpretationItems.Insert(setter).Exec(ctx, r.db)
		if err != nil {
			r.logger.ErrorContext(ctx, "Repository: Failed to create item",
				slog.String("item_id", item.ID),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("failed to create item %s: %w", item.ID, err)
		}
	}

	r.logger.InfoContext(ctx, "Repository: CreateItems completed",
		slog.Int("count", len(items)),
	)
	return nil
}

// UpdateItem はアイテムを更新します
func (r *interpretationItemRepository) UpdateItem(ctx context.Context, item *entity.InterpretationItem) error {
	r.logger.InfoContext(ctx, "Repository: UpdateItem started",
		slog.String("item_id", item.ID),
	)

	item.UpdatedAt = time.Now()

	setter := &models.InterpretationItemSetter{
		Data:      omit.From(types.JSON[json.RawMessage]{Val: item.Data}),
		UpdatedAt: omit.From(item.UpdatedAt),
	}

	_, err := models.InterpretationItems.Update(
		setter.UpdateMod(),
		um.Where(models.InterpretationItems.Columns.ID.EQ(mysql.Arg(item.ID))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to update item",
			slog.String("item_id", item.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to update item: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: UpdateItem completed",
		slog.String("item_id", item.ID),
	)
	return nil
}

// ApproveItem はアイテムを承認します（ステータス更新 + resource_id設定）
func (r *interpretationItemRepository) ApproveItem(ctx context.Context, itemID string, resourceID string) error {
	r.logger.InfoContext(ctx, "Repository: ApproveItem started",
		slog.String("item_id", itemID),
		slog.String("resource_id", resourceID),
	)

	now := time.Now()
	setter := &models.InterpretationItemSetter{
		Status:     omit.From(string(entity.ItemStatusCreated)),
		ResourceID: omitnull.FromNull(null.From(resourceID)),
		ReviewedAt: omitnull.FromNull(null.From(now)),
		UpdatedAt:  omit.From(now),
	}

	_, err := models.InterpretationItems.Update(
		setter.UpdateMod(),
		um.Where(models.InterpretationItems.Columns.ID.EQ(mysql.Arg(itemID))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to approve item",
			slog.String("item_id", itemID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to approve item: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: ApproveItem completed",
		slog.String("item_id", itemID),
	)
	return nil
}

// ApproveItems は複数のアイテムを承認します（トランザクション内で実行されることを想定）
func (r *interpretationItemRepository) ApproveItems(ctx context.Context, approvals map[string]string) error {
	r.logger.InfoContext(ctx, "Repository: ApproveItems started",
		slog.Int("count", len(approvals)),
	)

	now := time.Now()
	for itemID, resourceID := range approvals {
		setter := &models.InterpretationItemSetter{
			Status:     omit.From(string(entity.ItemStatusCreated)),
			ResourceID: omitnull.FromNull(null.From(resourceID)),
			ReviewedAt: omitnull.FromNull(null.From(now)),
			UpdatedAt:  omit.From(now),
		}

		_, err := models.InterpretationItems.Update(
			setter.UpdateMod(),
			um.Where(models.InterpretationItems.Columns.ID.EQ(mysql.Arg(itemID))),
		).Exec(ctx, r.db)

		if err != nil {
			r.logger.ErrorContext(ctx, "Repository: Failed to approve item in batch",
				slog.String("item_id", itemID),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("failed to approve item %s: %w", itemID, err)
		}
	}

	r.logger.InfoContext(ctx, "Repository: ApproveItems completed",
		slog.Int("count", len(approvals)),
	)
	return nil
}

// toEntity はDB modelをentityに変換します
func (r *interpretationItemRepository) toEntity(dbItem *models.InterpretationItem) (*entity.InterpretationItem, error) {
	item := &entity.InterpretationItem{
		ID:               dbItem.ID,
		InterpretationID: dbItem.InterpretationID,
		ItemIndex:        int(dbItem.ItemIndex),
		ResourceType:     entity.ResourceType(dbItem.ResourceType),
		Status:           entity.ItemStatus(dbItem.Status),
		Data:             dbItem.Data.Val,
		OriginalData:     dbItem.OriginalData.Val,
		CreatedAt:        dbItem.CreatedAt,
		UpdatedAt:        dbItem.UpdatedAt,
	}

	if dbItem.ResourceID.IsValue() {
		item.ResourceID = dbItem.ResourceID.Ptr()
	}

	if dbItem.ReviewedAt.IsValue() {
		item.ReviewedAt = dbItem.ReviewedAt.Ptr()
	}

	return item, nil
}
