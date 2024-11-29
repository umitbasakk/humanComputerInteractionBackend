package Controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/middlewares"
	AIModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	UserModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"

	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

type AIController struct {
	aiService     interfaces.AIService
	appMiddleware *middlewares.AppMiddleware
}

func NewAIController(echo *echo.Echo, service interfaces.AIService, middleware *middlewares.AppMiddleware) {
	Controller := &AIController{
		aiService:     service,
		appMiddleware: middleware,
	}

	echo.POST("/ai/requestAI", Controller.RequestAI, Controller.appMiddleware.AuthenticationMiddleware)
	echo.GET("/ai/getAllRequests", Controller.GetAllRequests, Controller.appMiddleware.AuthenticationMiddleware)
}

func (Controller *AIController) RequestAI(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*UserModel.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}
	request := &AIModel.AIRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: "No Bind", ErrCode: model.ErrorVerifySystem})
	}
	return Controller.aiService.GetResult(ctx.Request().Context(), ctx, request, user)
}

func (Controller *AIController) GetAllRequests(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*UserModel.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
	}

	return Controller.aiService.GetAllRequests(ctx.Request().Context(), ctx, user)
}
