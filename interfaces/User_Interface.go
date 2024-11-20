package interfaces

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
)

type UserService interface {
	GetUserByToken(ctx echo.Context, token string) error
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, user *model.User) error
}

type UserDatalayer interface {
	GetUserByID(ctx context.Context, userID int16) *model.User
	Signup(ctx echo.Context, user *model.User) error
	Login(ctx echo.Context, username string) (*model.User, error)
	GetUserByUsername(ctx echo.Context, username string) (*model.User, error)
	GetUserByEmail(ctx echo.Context, email string) (*model.User, error)
	SaveTokenByUsername(ctx echo.Context, username string, token string) error
}
