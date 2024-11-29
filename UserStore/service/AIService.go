package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	AIModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	AiModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	AuthModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"

	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

type AIServiceImpl struct {
	aiDL interfaces.AIDataLayer
}

func NewAIServiceImpl(dataLayer interfaces.AIDataLayer) interfaces.AIService {
	return &AIServiceImpl{
		aiDL: dataLayer,
	}
}

func (AIService *AIServiceImpl) GetResult(context context.Context, ctx echo.Context, request *AIModel.AIRequest, user *AuthModel.User) error {
	aiData := &AiModel.AIData{}

	aiData.UserId = strconv.Itoa(user.Id)
	aiData.StartedDate = request.StartedDate
	aiData.EndDate = request.EndDate
	aiData.HashTag = request.HashTag
	aiData.Category, _ = strconv.Atoi(request.Category)
	aiData.QuantityLimit, _ = strconv.Atoi(request.QuantityLimit)
	aiData.RequestStatus = 1
	/*
			url := "http://127.0.0.1:5000/endpoint"
			jsonVal, err := json.Marshal(request)
			if err != nil {
				fmt.Println("Hata", err)
			}
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))
			if err != nil {
				fmt.Println("Hata", err)
			}

		defer resp.Body.Close()
	*/
	tx, err := AIService.aiDL.GetTransaction(context)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.ErrorAI, ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	err = AIService.aiDL.SaveAiRequest(tx, ctx, aiData)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.ErrorAI, ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	err = AIService.aiDL.CommitTransaction(tx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.ErrorAI, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: "Devam ke", ErrCode: model.ErrorLoginSystem})
}

func (AIService *AIServiceImpl) GetAllRequests(context context.Context, ctx echo.Context, user *AuthModel.User) error {
	tx, err := AIService.aiDL.GetTransaction(context)
	if err != nil {
		return err
	}
	result, err := AIService.aiDL.GetRequestOfUser(tx, ctx, strconv.Itoa(user.Id))
	if err != nil {
		return err
	}
	if err := AIService.aiDL.CommitTransaction(tx); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: "", ErrCode: model.NoError, Data: result})
}
