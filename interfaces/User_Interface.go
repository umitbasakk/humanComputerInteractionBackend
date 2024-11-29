package interfaces

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
)

type UserService interface {
	Signup(context context.Context, ctx echo.Context, user *model.User) error
	Login(context context.Context, ctx echo.Context, user *model.User) error
	VerifyCode(context context.Context, ctx echo.Context, verify *model.VerifyRequest, user *model.User) error
	ResendCode(context context.Context, ctx echo.Context, user *model.User) error
	SendSms(ctx echo.Context, phone string, code string) error
	ChangePassword(ctx echo.Context, changePassword *model.PasswordRequest, user *model.User) error
	UpdateProfile(context context.Context, ctx echo.Context, profile *model.UpdateProfileRequest, user *model.User) error
}

type UserDatalayer interface {
	GetUserByID(ctx context.Context, userID int16) *model.User
	Signup(tx *sql.Tx, ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, username string) (*model.User, error)
	GetUserByPhone(tx *sql.Tx, ctx echo.Context, phone string) error
	IsThereEqualUsername(tx *sql.Tx, ctx echo.Context, username string) error
	IsThereEqualEmail(tx *sql.Tx, ctx echo.Context, email string) error
	GetUserUsername(tx *sql.Tx, ctx echo.Context, username string) (*model.User, error)
	GetUserEmail(tx *sql.Tx, ctx echo.Context, email string) (*model.User, error)
	SaveTokenByUsername(tx *sql.Tx, ctx echo.Context, username string, token string) error
	GetVerifyCode(tx *sql.Tx, ctx echo.Context, user_id int) (*model.Verify, error)
	CreateVerifyCode(tx *sql.Tx, ctx echo.Context, verify *model.Verify, user_id int) error
	VerifyCode(tx *sql.Tx, ctx echo.Context, user_id int) error
	UpdateVerifyCode(tx *sql.Tx, ctx echo.Context, user_id int, vCode string) error
	ChangePassword(ctx echo.Context, username string, password string) error
	UpdateProfile(tx *sql.Tx, ctx echo.Context, profile *model.UpdateProfileRequest, username string) error
	GetTransaction(ctx context.Context) (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
}
