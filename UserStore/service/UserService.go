package service

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/helpers"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userDL       interfaces.UserDatalayer
	twilioClient *twilio.RestClient
}

func NewUserServiceImpl(UserDatalayer interfaces.UserDatalayer, client *twilio.RestClient) interfaces.UserService {
	return &UserServiceImpl{userDL: UserDatalayer, twilioClient: client}
}

func (service *UserServiceImpl) Login(ctx echo.Context, user *model.User) error {

	result, err := service.userDL.Login(ctx, user.Username)
	if err != nil {
		return ctx.String(http.StatusOK, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidPassword, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	token, err := helpers.CreateJWTToken(user.Username)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	err = service.userDL.SaveTokenByUsername(ctx, token, user.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	verify, err := service.userDL.GetVerifyCode(ctx, user.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	userResponse := &model.UserResponse{
		Name:       result.Name,
		Username:   result.Username,
		Email:      result.Email,
		Phone:      result.Phone,
		Token:      token,
		Created_at: result.Created_at,
		Updated_at: result.Updated_at,
	}
	if verify.VerifyStatus != 1 {
		return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.MustbeVerified, ErrCode: model.MustbeVerified, Data: userResponse})
	}

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessLogin, ErrCode: model.NoError, Data: result})
}
func (service *UserServiceImpl) VerifyCode(ctx echo.Context, verify *model.VerifyRequest) error {

	result, err := service.userDL.GetVerifyCode(ctx, verify.Username)
	if result.VerifyStatus != 0 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsVerify, ErrCode: model.ErrorVerifySystem})
	}
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	if result.VerifyCode != verify.VerifyCode {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidVerifyCode, ErrCode: model.ErrorVerifySystem})
	}
	err = service.userDL.VerifyCode(ctx, verify.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessVerifyCode, ErrCode: model.NoError})
}

func (service *UserServiceImpl) ResendCode(ctx echo.Context, resendCodeRequest *model.ResendCodeRequest) error {
	result, err := service.userDL.GetVerifyCode(ctx, resendCodeRequest.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}

	if result.VerifyStatus != 0 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsVerify, ErrCode: model.ErrorVerifySystem})
	}
	timeDifference := time.Now().Sub(result.Updated_at).Seconds()
	if timeDifference < 120 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: fmt.Sprintf("you have to wait %v seconds for the next code.", math.Round(120-timeDifference)), ErrCode: model.ErrorVerifySystem})
	}

	vCode := strconv.Itoa(helpers.GetVerifyCode())
	err = service.userDL.UpdateVerifyCode(ctx, result.Username, vCode)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	err = service.SendSms(ctx, resendCodeRequest.Phone, vCode)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.SmsFailed, ErrCode: model.ErrorVerifySystem})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessResendCode, ErrCode: model.NoError})
}
func (service *UserServiceImpl) Signup(ctx echo.Context, user *model.User) error {

	user.Name = strings.ReplaceAll(user.Name, " ", "")
	user.Username = strings.ReplaceAll(user.Username, " ", "")
	user.Email = strings.ReplaceAll(user.Email, " ", "")

	if len(user.Username) < 5 || len(user.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientNameAndSurname, ErrCode: model.ErrorRegisterSystem})
	}

	if err := helpers.ValidEmail(user.Email); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidMail, ErrCode: model.ErrorRegisterSystem})
	}

	if len(user.Password) < 7 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientPassword, ErrCode: model.ErrorRegisterSystem})
	}

	if us, _ := service.userDL.GetUserByUsername(ctx, user.Username); us != nil && us.Username == user.Username {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsUsername, ErrCode: model.ErrorRegisterSystem})
	}
	if us, _ := service.userDL.GetUserByEmail(ctx, user.Email); us != nil && us.Email == user.Email {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsEmail, ErrCode: model.ErrorRegisterSystem})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorRegisterSystem})
	}

	user.Password = string(hash)
	err = service.userDL.Signup(ctx, user)
	if err != nil {
		return err
	}

	vCode := strconv.Itoa(helpers.GetVerifyCode())
	verify := &model.Verify{Username: user.Username, VerifyCode: vCode, VerifyStatus: 0}

	err = service.userDL.CreateVerifyCode(ctx, verify)
	if err != nil {
		return err
	}

	err = service.SendSms(ctx, user.Phone, verify.VerifyCode)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.ErrorVerified, ErrCode: model.ErrorVerifySystem, Data: nil})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessfullyRegistered, ErrCode: model.NoError})
}

func (service *UserServiceImpl) SendSms(ctx echo.Context, phone string, code string) error {

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom("+15109074928")
	params.SetBody(code)

	_, err := service.twilioClient.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	return nil
}
