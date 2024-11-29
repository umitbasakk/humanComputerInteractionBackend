package Controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/middlewares"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

type UserController struct {
	userService   interfaces.UserService
	appMiddleware *middlewares.AppMiddleware
}

func NewUserController(echoCtx *echo.Echo, userServiceObject interfaces.UserService, appMiddleware *middlewares.AppMiddleware) {

	userControllerObject := &UserController{
		userService:   userServiceObject,
		appMiddleware: appMiddleware,
	}
	echoCtx.POST("/register", userControllerObject.Signup)
	echoCtx.POST("/login", userControllerObject.Login)
	echoCtx.POST("/verify", userControllerObject.Verify, userControllerObject.appMiddleware.AuthenticationMiddleware)
	echoCtx.POST("/resendCode", userControllerObject.ResendCode, userControllerObject.appMiddleware.AuthenticationMiddleware)
	echoCtx.POST("/changePassword", userControllerObject.ChangePassword, userControllerObject.appMiddleware.AuthenticationMiddleware)
	echoCtx.POST("/updateProfile", userControllerObject.UpdateProfile, userControllerObject.appMiddleware.AuthenticationMiddleware)
	aasd, _ := bcrypt.GenerateFromPassword([]byte("135980Aa"), 10)
	log.Println(string(aasd))
}

func (userController *UserController) Signup(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}

	return userController.userService.Signup(ec.Request().Context(), ec, userM)
}

func (userController *UserController) Login(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}
	return userController.userService.Login(ec.Request().Context(), ec, userM)
}

func (userController *UserController) Verify(ec echo.Context) error {
	user, ok := ec.Get("user").(*model.User)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}
	verifyRequest := &model.VerifyRequest{}
	if err := ec.Bind(verifyRequest); err != nil {
		return err
	}
	return userController.userService.VerifyCode(ec.Request().Context(), ec, verifyRequest, user)
}

func (userController *UserController) ResendCode(ec echo.Context) error {
	user, ok := ec.Get("user").(*model.User)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}
	return userController.userService.ResendCode(ec.Request().Context(), ec, user)
}

func (userController *UserController) ChangePassword(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}
	passwordRq := &model.PasswordRequest{}
	if err := c.Bind(passwordRq); err != nil {
		return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	return userController.userService.ChangePassword(c, passwordRq, user)
}

func (userController *UserController) UpdateProfile(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}
	updateProfileRq := &model.UpdateProfileRequest{}
	if err := c.Bind(updateProfileRq); err != nil {
		return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	return userController.userService.UpdateProfile(c.Request().Context(), c, updateProfileRq, user)
}
