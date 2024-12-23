package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

var timeout = 300 * time.Second

func (AIService *AIServiceImpl) GetResult(context context.Context, ctx echo.Context, request *AIModel.AIRequest, user *model.User) error {

	log.Println("Coming Request")
	aiData := &AiModel.AIData{}
	aiData.UserId = strconv.Itoa(user.Id)
	aiData.StartedDate = request.StartedDate
	aiData.EndDate = request.EndDate
	aiData.HashTag = request.HashTag
	aiData.Category, _ = strconv.Atoi(request.Category)
	aiData.QuantityLimit, _ = strconv.Atoi(request.QuantityLimit)
	aiData.RequestStatus = 1

	url := fmt.Sprintf("http://%s:5000/getValue", os.Getenv("PYTHON_HOST"))

	jsonVal, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Hata", err)
	}

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:      10,
			IdleConnTimeout:   300 * time.Second,
			DisableKeepAlives: false,
		},
		Timeout: timeout,
	}

	log.Println("Posting Request")

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonVal))
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	log.Println("Before Defer")
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	log.Println("After Defer")
	log.Println(string(s))
	response := &AIModel.AIResponse{}
	err = json.Unmarshal(s, response)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	log.Println(string(s))

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.Successful, ErrCode: model.ErrorLoginSystem, Data: response})
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
