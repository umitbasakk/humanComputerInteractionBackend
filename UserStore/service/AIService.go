package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	AIModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	AuthModel "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
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

func (AIService *AIServiceImpl) GetResult(ctx echo.Context, request *AIModel.AIRequest, user *AuthModel.User) error {
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
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: "Devam ke", ErrCode: model.ErrorLoginSystem})
}
