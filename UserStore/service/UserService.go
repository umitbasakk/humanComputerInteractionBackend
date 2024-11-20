package service

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/helpers"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userDL interfaces.UserDatalayer
}

func NewUserServiceImpl(UserDatalayer interfaces.UserDatalayer) interfaces.UserService {
	return &UserServiceImpl{userDL: UserDatalayer}
}

func (service *UserServiceImpl) GetUserByToken(ctx echo.Context, token string) error {
	return ctx.String(http.StatusOK, token)
}

func (service *UserServiceImpl) Login(ctx echo.Context, user *model.User) error {

	result, err := service.userDL.Login(ctx, user.Username)
	if err != nil {
		return ctx.String(http.StatusOK, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Username or Password incorrect")
	}

	token, err := helpers.CreateJWTToken(user.Username)

	if err != nil {
		return err
	}

	err = service.userDL.SaveTokenByUsername(ctx, token, user.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: "success", Data: result})

}
func (service *UserServiceImpl) Signup(ctx echo.Context, user *model.User) error {

	user.Name = strings.ReplaceAll(user.Name, " ", "")
	user.Username = strings.ReplaceAll(user.Username, " ", "")
	user.Email = strings.ReplaceAll(user.Email, " ", "")

	if len(user.Username) < 5 || len(user.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientNameAndSurname})
	}

	if err := helpers.ValidEmail(user.Email); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidMail})
	}

	if len(user.Password) < 7 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientPassword})
	}

	if us, _ := service.userDL.GetUserByUsername(ctx, user.Username); us != nil && us.Username == user.Username {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsUsername})
	}
	if us, _ := service.userDL.GetUserByEmail(ctx, user.Email); us != nil && us.Email == user.Email {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsEmail})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return ctx.String(http.StatusBadRequest, "failed to hash password")
	}

	user.Password = string(hash)
	err = service.userDL.Signup(ctx, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessfullyRegistered})
}
