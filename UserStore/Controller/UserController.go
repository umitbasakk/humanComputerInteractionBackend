package Controller

import (
	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(echoCtx *echo.Echo, userServiceObject interfaces.UserService) {

	userControllerObject := &UserController{
		userService: userServiceObject,
	}

	echoCtx.GET("/getUser/:token", userControllerObject.GetUser)
	echoCtx.POST("/signup", userControllerObject.Signup)
	echoCtx.POST("/login", userControllerObject.Login)

}

func (user *UserController) GetUser(ec echo.Context) error {
	token := ec.Param("token")
	return user.userService.GetUserByToken(ec, token)
}

func (user *UserController) Signup(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}

	return user.userService.Signup(ec, userM)
}

func (user *UserController) Login(ec echo.Context) error {
	userM := &model.User{}
	if err := ec.Bind(userM); err != nil {
		return err
	}
	return user.userService.Login(ec, userM)
}
