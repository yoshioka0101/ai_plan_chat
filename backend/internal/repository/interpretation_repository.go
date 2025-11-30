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
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

type interpretationRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

// NewInterpretationRepository は新しいInterpretationRepositoryを生成します
func NewInterpretationRepository(db *sql.DB, logger *slog.Logger) interfaces.InterpretationRepository {
	return NewInterpretationRepositoryWithExecutor(bob.NewDB(db), logger)
}

// NewInterpretationRepositoryWithExecutor は既存のexecutorを使ってInterpretationRepositoryを生成します
func NewInterpretationRepositoryWithExecutor(exec bob.Executor, logger *slog.Logger) interfaces.InterpretationRepository {
	return &interpretationRepository{
		db:     exec,
		logger: logger,
	}
}

// CreateInterpretation はAI解釈結果を保存します
func (r *interpretationRepository) CreateInterpretation(ctx context.Context, interpretation *entity.AIInterpretation) error {
	r.logger.InfoContext(ctx, "Repository: CreateInterpretation started",
		slog.String("interpretation_id", interpretation.ID),
		slog.String("user_id", interpretation.UserID),
	)

	// InterpretationResultをJSONに変換
	structuredResult, err := json.Marshal(map[string]interface{}{
		"type":        interpretation.Result.Type,
		"title":       interpretation.Result.Title,
		"description": interpretation.Result.Description,
		"metadata":    interpretation.Result.Metadata,
	})
	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to marshal structured result",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to marshal structured result: %w", err)
	}

	// 現在時刻を設定
	now := time.Now()
	interpretation.CreatedAt = now
	interpretation.UpdatedAt = now

	// AIトークン数をnull.Valに変換
	var aiPromptTokens null.Val[int32]
	if interpretation.AIPromptTokens != nil {
		aiPromptTokens = null.From(int32(*interpretation.AIPromptTokens))
	}

	var aiCompletionTokens null.Val[int32]
	if interpretation.AICompletionTokens != nil {
		aiCompletionTokens = null.From(int32(*interpretation.AICompletionTokens))
	}

	// OriginalResultをnull.Valに変換
	var originalResult null.Val[types.JSON[json.RawMessage]]
	if len(interpretation.OriginalResult) > 0 {
		originalResult = null.From(types.NewJSON(json.RawMessage(interpretation.OriginalResult)))
	}

	_, err = models.AiInterpretations.Insert(
		&models.AiInterpretationSetter{
			ID:                 omit.From(interpretation.ID),
			UserID:             omit.From(interpretation.UserID),
			InputText:          omit.From(interpretation.InputText),
			StructuredResult:   omit.From(types.NewJSON(json.RawMessage(structuredResult))),
			OriginalResult:     omitnull.FromNull(originalResult),
			AiModel:            omit.From(interpretation.AIModel),
			AiPromptTokens:     omitnull.FromNull(aiPromptTokens),
			AiCompletionTokens: omitnull.FromNull(aiCompletionTokens),
			CreatedAt:          omit.From(interpretation.CreatedAt),
		},
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to create interpretation",
			slog.String("interpretation_id", interpretation.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to create interpretation: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: CreateInterpretation completed",
		slog.String("interpretation_id", interpretation.ID),
	)
	return nil
}

// GetInterpretationByID はIDで解釈結果を取得します
func (r *interpretationRepository) GetInterpretationByID(ctx context.Context, id string) (*entity.AIInterpretation, error) {
	r.logger.InfoContext(ctx, "Repository: GetInterpretationByID started",
		slog.String("interpretation_id", id),
	)

	aiInterpretation, err := models.AiInterpretations.Query(
		sm.Where(models.AiInterpretations.Columns.ID.EQ(mysql.Arg(id))),
	).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "Repository: Interpretation not found",
				slog.String("interpretation_id", id),
			)
			return nil, fmt.Errorf("interpretation not found: %s", id)
		}
		r.logger.ErrorContext(ctx, "Repository: Failed to query interpretation",
			slog.String("interpretation_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find interpretation: %w", err)
	}

	// BOBモデルからEntityに変換
	entityInterpretation, err := r.toEntity(aiInterpretation)
	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to convert to entity",
			slog.String("interpretation_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to convert to entity: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetInterpretationByID completed",
		slog.String("interpretation_id", id),
	)
	return entityInterpretation, nil
}

// GetInterpretationsByUserID はユーザーIDで解釈結果のリストを取得します
func (r *interpretationRepository) GetInterpretationsByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AIInterpretation, error) {
	r.logger.InfoContext(ctx, "Repository: GetInterpretationsByUserID started",
		slog.String("user_id", userID),
		slog.Int("limit", limit),
		slog.Int("offset", offset),
	)

	aiInterpretations, err := models.AiInterpretations.Query(
		sm.Where(models.AiInterpretations.Columns.UserID.EQ(mysql.Arg(userID))),
		sm.OrderBy(mysql.Raw("created_at DESC")),
		sm.Limit(int64(limit)),
		sm.Offset(int64(offset)),
	).All(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to query interpretations",
			slog.String("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get interpretations: %w", err)
	}

	// BOBモデルからEntityに変換
	var result []*entity.AIInterpretation
	for _, ai := range aiInterpretations {
		entityInterpretation, err := r.toEntity(ai)
		if err != nil {
			r.logger.ErrorContext(ctx, "Repository: Failed to convert to entity",
				slog.String("interpretation_id", ai.ID),
				slog.String("error", err.Error()),
			)
			continue // スキップして次へ
		}
		result = append(result, entityInterpretation)
	}

	r.logger.InfoContext(ctx, "Repository: GetInterpretationsByUserID completed",
		slog.String("user_id", userID),
		slog.Int("count", len(result)),
	)
	return result, nil
}

// toEntity はBOBモデルをEntityに変換します
func (r *interpretationRepository) toEntity(ai *models.AiInterpretation) (*entity.AIInterpretation, error) {
	// JSONからInterpretationResultに変換
	var structuredResult struct {
		Type        entity.InterpretationType     `json:"type"`
		Title       string                        `json:"title"`
		Description string                        `json:"description"`
		Metadata    entity.InterpretationMetadata `json:"metadata"`
	}

	structuredResultBytes, _ := ai.StructuredResult.MarshalJSON()
	if err := json.Unmarshal(structuredResultBytes, &structuredResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal structured result: %w", err)
	}

	// AIトークン数をポインタに変換
	var aiPromptTokens *int
	promptTokensVal, promptTokensNull := ai.AiPromptTokens.Get()
	if !promptTokensNull {
		val := int(promptTokensVal)
		aiPromptTokens = &val
	}

	var aiCompletionTokens *int
	completionTokensVal, completionTokensNull := ai.AiCompletionTokens.Get()
	if !completionTokensNull {
		val := int(completionTokensVal)
		aiCompletionTokens = &val
	}

	var originalResult []byte
	if val, ok := ai.OriginalResult.Get(); ok {
		raw, _ := val.MarshalJSON()
		originalResult = raw
	}

	return &entity.AIInterpretation{
		ID:        ai.ID,
		UserID:    ai.UserID,
		InputText: ai.InputText,
		Result: entity.InterpretationResult{
			Type:        structuredResult.Type,
			Title:       structuredResult.Title,
			Description: structuredResult.Description,
			Metadata:    structuredResult.Metadata,
		},
		OriginalResult:     originalResult,
		AIModel:            ai.AiModel,
		AIPromptTokens:     aiPromptTokens,
		AICompletionTokens: aiCompletionTokens,
		CreatedAt:          ai.CreatedAt,
		UpdatedAt:          ai.CreatedAt, // created_atのみなのでupdated_atも同じ値
	}, nil
}
