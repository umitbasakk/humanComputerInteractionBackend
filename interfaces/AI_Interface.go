package interfaces

import (
	"github.com/labstack/echo/v4"
	AIModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	AuthModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
)

type AIService interface {
	GetResult(ctx echo.Context, request *AIModel.AIRequest, user *AuthModel.User) error
}

type AIDataLayer interface {
}
