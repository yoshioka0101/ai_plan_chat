package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/database"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
	"github.com/yoshioka0101/ai_plan_chat/internal/repository"
)

type interpretationItemUseCase struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewInterpretationItemUseCase は新しいInterpretationItemUseCaseを生成します
func NewInterpretationItemUseCase(db *sql.DB, logger *slog.Logger) interfaces.InterpretationItemUseCase {
	return &interpretationItemUseCase{
		db:     db,
		logger: logger,
	}
}

// GetItems は指定したAI解釈IDに紐づくアイテム一覧を取得します
func (uc *interpretationItemUseCase) GetItems(ctx context.Context, interpretationID string) ([]*entity.InterpretationItem, error) {
	uc.logger.InfoContext(ctx, "UseCase: GetItems started",
		slog.String("interpretation_id", interpretationID),
	)

	itemRepo := repository.NewInterpretationItemRepository(bob.NewDB(uc.db), uc.logger)
	items, err := itemRepo.GetItemsByInterpretationID(ctx, interpretationID)
	if err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to get items",
			slog.String("interpretation_id", interpretationID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "UseCase: GetItems completed",
		slog.String("interpretation_id", interpretationID),
		slog.Int("count", len(items)),
	)
	return items, nil
}

// GetItem はIDでアイテムを取得します
func (uc *interpretationItemUseCase) GetItem(ctx context.Context, itemID string) (*entity.InterpretationItem, error) {
	uc.logger.InfoContext(ctx, "UseCase: GetItem started",
		slog.String("item_id", itemID),
	)

	itemRepo := repository.NewInterpretationItemRepository(bob.NewDB(uc.db), uc.logger)
	item, err := itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to get item",
			slog.String("item_id", itemID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "UseCase: GetItem completed",
		slog.String("item_id", itemID),
	)
	return item, nil
}

// UpdateItem はアイテムのdataフィールドを更新します
func (uc *interpretationItemUseCase) UpdateItem(ctx context.Context, itemID string, data []byte) (*entity.InterpretationItem, error) {
	uc.logger.InfoContext(ctx, "UseCase: UpdateItem started",
		slog.String("item_id", itemID),
	)

	itemRepo := repository.NewInterpretationItemRepository(bob.NewDB(uc.db), uc.logger)

	// 既存アイテムを取得
	item, err := itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to get item",
			slog.String("item_id", itemID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// pending状態のみ編集可能
	if item.Status != entity.ItemStatusPending {
		uc.logger.WarnContext(ctx, "UseCase: Cannot update non-pending item",
			slog.String("item_id", itemID),
			slog.String("status", string(item.Status)),
		)
		return nil, fmt.Errorf("cannot update item with status: %s", item.Status)
	}

	// データ更新
	item.Data = data

	// 更新実行
	if err := itemRepo.UpdateItem(ctx, item); err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to update item",
			slog.String("item_id", itemID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "UseCase: UpdateItem completed",
		slog.String("item_id", itemID),
	)
	return item, nil
}

// ApproveItem はアイテムを承認してリソース（タスク等）を作成します
func (uc *interpretationItemUseCase) ApproveItem(ctx context.Context, itemID string) (string, error) {
	uc.logger.InfoContext(ctx, "UseCase: ApproveItem started",
		slog.String("item_id", itemID),
	)

	var resourceID string

	err := database.WithTransaction(ctx, uc.db, func(tx bob.Executor) error {
		itemRepo := repository.NewInterpretationItemRepository(tx, uc.logger)

		// アイテム取得
		item, err := itemRepo.GetItemByID(ctx, itemID)
		if err != nil {
			return fmt.Errorf("failed to get item: %w", err)
		}

		// pending状態のみ承認可能
		if item.Status != entity.ItemStatusPending {
			return fmt.Errorf("cannot approve item with status: %s", item.Status)
		}

		// リソースタイプに応じてリソース作成
		var createdResourceID string
		switch item.ResourceType {
		case entity.ResourceTypeTask:
			createdResourceID, err = uc.createTaskFromItem(ctx, tx, item)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
		default:
			return fmt.Errorf("unsupported resource type: %s", item.ResourceType)
		}

		// アイテムを承認済みに更新
		if err := itemRepo.ApproveItem(ctx, itemID, createdResourceID); err != nil {
			return fmt.Errorf("failed to approve item: %w", err)
		}

		resourceID = createdResourceID
		return nil
	})

	if err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to approve item",
			slog.String("item_id", itemID),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	uc.logger.InfoContext(ctx, "UseCase: ApproveItem completed",
		slog.String("item_id", itemID),
		slog.String("resource_id", resourceID),
	)
	return resourceID, nil
}

// ApproveMultipleItems は複数のアイテムを一括承認します（トランザクション）
func (uc *interpretationItemUseCase) ApproveMultipleItems(ctx context.Context, itemIDs []string) (map[string]string, error) {
	uc.logger.InfoContext(ctx, "UseCase: ApproveMultipleItems started",
		slog.Int("count", len(itemIDs)),
	)

	resourceIDs := make(map[string]string)

	err := database.WithTransaction(ctx, uc.db, func(tx bob.Executor) error {
		itemRepo := repository.NewInterpretationItemRepository(tx, uc.logger)

		for _, itemID := range itemIDs {
			// アイテム取得
			item, err := itemRepo.GetItemByID(ctx, itemID)
			if err != nil {
				return fmt.Errorf("failed to get item %s: %w", itemID, err)
			}

			// pending状態のみ承認可能
			if item.Status != entity.ItemStatusPending {
				return fmt.Errorf("cannot approve item %s with status: %s", itemID, item.Status)
			}

			// リソースタイプに応じてリソース作成
			var createdResourceID string
			switch item.ResourceType {
			case entity.ResourceTypeTask:
				createdResourceID, err = uc.createTaskFromItem(ctx, tx, item)
				if err != nil {
					return fmt.Errorf("failed to create task for item %s: %w", itemID, err)
				}
			default:
				return fmt.Errorf("unsupported resource type for item %s: %s", itemID, item.ResourceType)
			}

			// アイテムを承認済みに更新
			if err := itemRepo.ApproveItem(ctx, itemID, createdResourceID); err != nil {
				return fmt.Errorf("failed to approve item %s: %w", itemID, err)
			}

			resourceIDs[itemID] = createdResourceID
		}

		return nil
	})

	if err != nil {
		uc.logger.ErrorContext(ctx, "UseCase: Failed to approve multiple items",
			slog.Int("count", len(itemIDs)),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "UseCase: ApproveMultipleItems completed",
		slog.Int("count", len(itemIDs)),
	)
	return resourceIDs, nil
}

// createTaskFromItem はアイテムからタスクを作成します（トランザクション内で実行）
func (uc *interpretationItemUseCase) createTaskFromItem(ctx context.Context, tx bob.Executor, item *entity.InterpretationItem) (string, error) {
	// トランザクション内で動作するRepositoryを作成
	taskRepo := repository.NewTaskRepositoryWithExecutor(tx, uc.logger)

	// JSONデータをTaskDataにパース
	var taskData entity.TaskData
	if err := json.Unmarshal(item.Data, &taskData); err != nil {
		return "", fmt.Errorf("failed to parse task data: %w", err)
	}

	// contextからuserIDを取得
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user_id not found in context")
	}

	// タスク作成
	task := &models.Task{
		ID:                 uuid.New().String(),
		UserID:             userID,
		Title:              taskData.Title,
		Description:        null.FromPtr(taskData.Description),
		DueAt:              null.FromPtr(taskData.DueAt),
		Status:             resolveTaskStatus(taskData.Status),
		Source:             "ai",
		AiInterpretationID: null.From(item.InterpretationID),
	}

	if err := taskRepo.CreateTask(ctx, task); err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	return task.ID, nil
}

// resolveTaskStatus validates and returns a task status, defaulting to "todo" for invalid or empty values.
func resolveTaskStatus(statusPtr *string) string {
	if statusPtr == nil || *statusPtr == "" {
		return "todo"
	}

	switch *statusPtr {
	case "todo", "in_progress", "done":
		return *statusPtr
	default:
		return "todo"
	}
}
