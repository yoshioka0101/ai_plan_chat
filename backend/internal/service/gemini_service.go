package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
	"google.golang.org/api/option"
)

//go:embed prompts/*.tmpl
var promptFiles embed.FS

var promptTemplate *template.Template

// GeminiService はGemini APIとのやり取りを担当するサービス
type GeminiService struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	modelName string
}

// init はプロンプトテンプレートを初期化します
func init() {
	var err error
	promptTemplate, err = template.ParseFS(promptFiles, "prompts/*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("failed to parse prompt templates: %v", err))
	}
}

// NewGeminiService は新しいGeminiServiceを作成します
// modelNameは必須パラメータです。デフォルト値の設定はconfig層で行います。
func NewGeminiService(apiKey string, modelName string) (*GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	if modelName == "" {
		return nil, fmt.Errorf("Gemini model name is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.7)
	model.ResponseMIMEType = "application/json"

	return &GeminiService{
		client:    client,
		model:     model,
		modelName: modelName,
	}, nil
}

// InterpretInputResult はAI解釈の結果と生JSONを含む
type InterpretInputResult struct {
	Result       *entity.InterpretationResult
	OriginalJSON []byte
}

// InterpretInput はユーザーの入力を解析します
func (s *GeminiService) InterpretInput(ctx context.Context, inputText string) (*InterpretInputResult, error) {
	prompt := buildPrompt(inputText)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	// レスポンスをJSON形式でパース
	var rawResult struct {
		Type        string                 `json:"type"`
		Title       string                 `json:"title"`
		Description string                 `json:"description,omitempty"`
		Metadata    map[string]interface{} `json:"metadata,omitempty"`
	}

	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	originalJSON := []byte(responseText)

	if err := json.Unmarshal(originalJSON, &rawResult); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w, response: %s", err, responseText)
	}

	// entity型に変換
	result := &entity.InterpretationResult{
		Type:        entity.InterpretationType(rawResult.Type),
		Title:       rawResult.Title,
		Description: rawResult.Description,
		Metadata:    convertToInterpretationMetadata(rawResult.Metadata),
	}

	if result.Type == "" {
		result.Type = entity.InterpretationTypeTodo
	}

	return &InterpretInputResult{
		Result:       result,
		OriginalJSON: originalJSON,
	}, nil
}

// ModelName は使用中のモデル名を返します
func (s *GeminiService) ModelName() string {
	return s.modelName
}

// Close はGeminiクライアントをクローズします
func (s *GeminiService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// convertToInterpretationMetadata はmap[string]interface{}をInterpretationMetadataに変換（Todoのみ）
func convertToInterpretationMetadata(raw map[string]interface{}) entity.InterpretationMetadata {
	metadata := entity.InterpretationMetadata{
		Extra: make(map[string]interface{}),
	}

	if raw == nil {
		return metadata
	}

	// タグの変換
	if tags, ok := raw["tags"].([]interface{}); ok {
		tagStrings := make([]string, 0, len(tags))
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				tagStrings = append(tagStrings, tagStr)
			}
		}
		if len(tagStrings) > 0 {
			metadata.Tags = tagStrings
		}
	}

	// Todo関連フィールドの変換
	if priority, ok := raw["priority"].(string); ok {
		metadata.Priority = &priority
	}

	// 期限の変換
	if deadlineStr, ok := raw["deadline"].(string); ok {
		if t, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
			metadata.Deadline = &t
		}
	}

	// その他のフィールドはExtraに格納
	for key, value := range raw {
		switch key {
		case "tags", "priority", "deadline":
			// 既に処理済みまたは別途処理
		default:
			metadata.Extra[key] = value
		}
	}

	return metadata
}

// buildPrompt は解析用のプロンプトを構築します
func buildPrompt(inputText string) string {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"Input": inputText,
	}

	err := promptTemplate.ExecuteTemplate(&buf, "interpretation.tmpl", data)
	if err != nil {
		// フォールバック: テンプレートエラーの場合はシンプルなプロンプトを返す
		return fmt.Sprintf("ユーザーの入力を解析してください: %s", inputText)
	}

	return buf.String()
}
