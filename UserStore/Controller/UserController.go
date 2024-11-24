package Controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/middlewares"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
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
	echoCtx.POST("/verify", userControllerObject.Verify)
	echoCtx.POST("/resendCode", userControllerObject.ResendCode)
	echoCtx.POST("/changePassword", userControllerObject.ChangePassword)
	echoCtx.POST("/test", userControllerObject.Test, userControllerObject.appMiddleware.AuthenticationMiddleware)

}

func (userController *UserController) Signup(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}

	return userController.userService.Signup(ec, userM)
}

func (userController *UserController) Login(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}
	return userController.userService.Login(ec, userM)
}

func (userController *UserController) Verify(ec echo.Context) error {

	verifyRequest := &model.VerifyRequest{}
	if err := ec.Bind(verifyRequest); err != nil {
		return err
	}
	return userController.userService.VerifyCode(ec, verifyRequest)
}

func (userController *UserController) ResendCode(ec echo.Context) error {

	resendCodeRequest := &model.ResendCodeRequest{}
	if err := ec.Bind(resendCodeRequest); err != nil {
		return err
	}
	return userController.userService.ResendCode(ec, resendCodeRequest)
}

func (userController *UserController) ChangePassword(ec echo.Context) error {

	passwordRequest := &model.PasswordRequest{}
	if err := ec.Bind(passwordRequest); err != nil {
		return err
	}
	return nil
}

func (userController *UserController) Test(c echo.Context) error {
	_, ok := c.Get("user").(*model.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	passwordRequest := &model.PasswordRequest{}
	if err := c.Bind(passwordRequest); err != nil {
		return err
	}
	return nil
}
