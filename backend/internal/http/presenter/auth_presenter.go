package presenter

import (
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/models"
)

// AuthPresenter は認証のレスポンス整形を担当します
type AuthPresenter struct{}

func NewAuthPresenter() *AuthPresenter {
	return &AuthPresenter{}
}

// GoogleCallback はGoogleOAuthコールバックのレスポンスを整形します
func (p *AuthPresenter) GoogleCallback(user *models.User, jwtToken string) api.AuthResponse {
	// 文字列IDをUUIDにパース
	id, err := uuid.Parse(user.ID)
	if err != nil {
		// DB整合性が保たれていれば発生しないはず
		id = uuid.Nil
	}

	response := api.AuthResponse{
		JwtToken: jwtToken,
		User: api.User{
			Id:        types.UUID(id),
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt,
		},
	}

	// Null可能なフィールドの処理
	if val, ok := user.Avatar.Get(); ok {
		response.User.Avatar = &val
	}

	return response
}
