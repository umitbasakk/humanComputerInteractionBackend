package interfaces

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
)

type UserService interface {
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, user *model.User) error
	VerifyCode(ctx echo.Context, verify *model.VerifyRequest, user *model.User) error
	ResendCode(ctx echo.Context, user *model.User) error
	SendSms(ctx echo.Context, phone string, code string) error
	ChangePassword(ctx echo.Context, changePassword *model.PasswordRequest, user *model.User) error
	UpdateProfile(ctx echo.Context, profile *model.UpdateProfileRequest, user *model.User) error
}

type UserDatalayer interface {
	GetUserByID(ctx context.Context, userID int16) *model.User
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, username string) (*model.User, error)
	GetUserByPhone(ctx echo.Context, phone string) error
	IsThereEqualUsername(ctx echo.Context, username string) error
	IsThereEqualEmail(ctx echo.Context, email string) error
	GetUserUsername(ctx echo.Context, username string) (*model.User, error)
	GetUserEmail(ctx echo.Context, email string) (*model.User, error)
	SaveTokenByUsername(ctx echo.Context, username string, token string) error
	GetVerifyCode(ctx echo.Context, user_id int) (*model.Verify, error)
	CreateVerifyCode(ctx echo.Context, verify *model.Verify, user_id int) error
	VerifyCode(ctx echo.Context, user_id int) error
	UpdateVerifyCode(ctx echo.Context, user_id int, vCode string) error
	ChangePassword(ctx echo.Context, username string, password string) error
	UpdateProfile(ctx echo.Context, profile *model.UpdateProfileRequest, username string) error
}
