package interfaces

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
)

type UserService interface {
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, user *model.User) error
	VerifyCode(ctx echo.Context, verify *model.VerifyRequest) error
	ResendCode(ctx echo.Context, resendCodeRequest *model.ResendCodeRequest) error
	SendSms(ctx echo.Context, phone string, code string) error
}

type UserDatalayer interface {
	GetUserByID(ctx context.Context, userID int16) *model.User
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, username string) (*model.User, error)
	GetUserByUsername(ctx echo.Context, username string) (*model.User, error)
	GetUserByEmail(ctx echo.Context, email string) (*model.User, error)
	SaveTokenByUsername(ctx echo.Context, username string, token string) error
	GetVerifyCode(ctx echo.Context, username string) (*model.Verify, error)
	CreateVerifyCode(ctx echo.Context, verify *model.Verify) error
	VerifyCode(ctx echo.Context, username string) error
	UpdateVerifyCode(ctx echo.Context, username string, vCode string) error
}
