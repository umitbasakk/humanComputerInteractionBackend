package interfaces

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	AIModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	AuthModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
)

type AIService interface {
	GetResult(context context.Context, ctx echo.Context, request *AIModel.AIRequest, user *AuthModel.User) error
	GetAllRequests(context context.Context, ctx echo.Context, user *AuthModel.User) error
}

type AIDataLayer interface {
	SaveAiRequest(tx *sql.Tx, ctx echo.Context, aiData *model.AIData) error
	GetRequestOfUser(tx *sql.Tx, ctx echo.Context, user_id string) ([]AIModel.AIData, error)
	GetTransaction(ctx context.Context) (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
}
