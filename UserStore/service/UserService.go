package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
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

func (service *UserServiceImpl) Login(context context.Context, ctx echo.Context, user *model.User) error {

	result, err := service.userDL.Login(ctx, user.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.NotCorrectPassword, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidPassword, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	token, err := helpers.CreateJWTToken(result.Username)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: "a", ErrCode: model.ErrorLoginSystem, Data: nil})
	}
	tx, err := service.userDL.GetTransaction(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	err = service.userDL.SaveTokenByUsername(tx, ctx, token, result.Username)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorLoginSystem, Data: nil})
	}

	verify, err := service.userDL.GetVerifyCode(tx, ctx, result.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: "c", ErrCode: model.ErrorLoginSystem, Data: nil})
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
	if err := service.userDL.CommitTransaction(tx); err != nil {
		return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.MustbeVerified, Data: nil})

	}

	if verify.VerifyStatus != 1 {
		return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.MustbeVerified, ErrCode: model.MustbeVerified, Data: userResponse})
	}

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessLogin, ErrCode: model.NoError, Data: userResponse})
}
func (service *UserServiceImpl) VerifyCode(context context.Context, ctx echo.Context, verify *model.VerifyRequest, user *model.User) error {

	tx, err := service.userDL.GetTransaction(context)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}

	result, err := service.userDL.GetVerifyCode(tx, ctx, user.Id)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}

	if result.VerifyStatus != 0 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsVerify, ErrCode: model.ErrorVerifySystem})
	}

	if result.VerifyCode != verify.VerifyCode {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidVerifyCode, ErrCode: model.ErrorVerifySystem})
	}
	err = service.userDL.VerifyCode(tx, ctx, user.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	user, err = service.userDL.GetUserUsername(tx, ctx, user.Username)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}

	if err := service.userDL.CommitTransaction(tx); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}

	userResponse := &model.UserResponse{
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Phone:      user.Phone,
		Token:      user.Token,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessVerifyCode, ErrCode: model.NoError, Data: userResponse})
}

func (service *UserServiceImpl) ResendCode(context context.Context, ctx echo.Context, user *model.User) error {
	tx, err := service.userDL.GetTransaction(context)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}
	result, err := service.userDL.GetVerifyCode(tx, ctx, user.Id)
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
	err = service.userDL.UpdateVerifyCode(tx, ctx, user.Id, vCode)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}

	err = service.SendSms(ctx, user.Phone, vCode)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.SmsFailed, ErrCode: model.ErrorVerifySystem})
	}

	if err := service.userDL.CommitTransaction(tx); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessResendCode, ErrCode: model.NoError})
}
func (service *UserServiceImpl) Signup(context context.Context, ctx echo.Context, user *model.User) error {

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
	tx, err := service.userDL.GetTransaction(context)

	if errUseUsername := service.userDL.IsThereEqualUsername(tx, ctx, user.Username); errUseUsername != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsUsername, ErrCode: model.ErrorRegisterSystem})
	}
	if errUseEmail := service.userDL.IsThereEqualEmail(tx, ctx, user.Email); errUseEmail != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsEmail, ErrCode: model.ErrorRegisterSystem})
	}

	if errUsePhone := service.userDL.GetUserByPhone(tx, ctx, user.Phone); errUsePhone != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyUsedPhone, ErrCode: model.ErrorRegisterSystem})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorRegisterSystem})
	}

	user.Password = string(hash)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorVerifySystem, Data: nil})

	}
	err = service.userDL.Signup(tx, ctx, user) //save
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorVerifySystem, Data: nil})
	}

	vCode := strconv.Itoa(helpers.GetVerifyCode())
	verify := &model.Verify{Username: user.Username, VerifyCode: vCode, VerifyStatus: 0}
	userGt, errGetUser := service.userDL.GetUserEmail(tx, ctx, user.Email)
	if errGetUser != nil {
		log.Println(errGetUser)
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.SmsFailed, ErrCode: model.ErrorVerifySystem, Data: nil})
	}
	err = service.userDL.CreateVerifyCode(tx, ctx, verify, userGt.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.SmsFailed, ErrCode: model.ErrorVerifySystem, Data: nil})
	}

	err = service.SendSms(ctx, user.Phone, verify.VerifyCode)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorVerifySystem, Data: nil})
	}
	service.userDL.CommitTransaction(tx)

	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.SuccessfullyRegistered, ErrCode: model.NoError})
}

func (service *UserServiceImpl) ChangePassword(ctx echo.Context, changePassword *model.PasswordRequest, user *model.User) error {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassword.CurrentPassword)); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.NotCorrectPassword, ErrCode: model.ErrorVerifySystem})
	}
	if len(changePassword.NewPassword) < 7 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientPassword, ErrCode: model.ErrorVerifySystem})
	}
	if changePassword.CurrentPassword == changePassword.NewPassword {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.EqualPassword, ErrCode: model.ErrorVerifySystem})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	errProcess := service.userDL.ChangePassword(ctx, user.Username, string(hashPassword))
	if errProcess != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.GlobalError, ErrCode: model.ErrorVerifySystem})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.ChangedPassword, ErrCode: model.ErrorVerifySystem})
}

func (service *UserServiceImpl) UpdateProfile(context context.Context, ctx echo.Context, profile *model.UpdateProfileRequest, user *model.User) error {

	if len(profile.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientName, ErrCode: model.ErrorVerifySystem})
	}
	if err := helpers.ValidEmail(profile.Email); err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InvalidMail, ErrCode: model.ErrorVerifySystem})
	}

	if len(profile.Username) < 5 {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.InsufficientUsername, ErrCode: model.ErrorVerifySystem})
	}

	tx, err := service.userDL.GetTransaction(context)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}

	if err := service.userDL.IsThereEqualUsername(tx, ctx, profile.Username); err != nil && user.Username != profile.Username {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsUsername, ErrCode: model.ErrorVerifySystem})
	}

	if err := service.userDL.IsThereEqualEmail(tx, ctx, profile.Email); err != nil && user.Email != profile.Email {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.AlreadyExistsEmail, ErrCode: model.ErrorVerifySystem})
	}
	err = service.userDL.UpdateProfile(tx, ctx, profile, user.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: err.Error(), ErrCode: model.ErrorVerifySystem})
	}
	err = service.userDL.CommitTransaction(tx)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, &model.MessageHandler{Message: constants.FailedTransaction, ErrCode: model.ErrorVerifySystem})
	}
	return ctx.JSON(http.StatusOK, &model.MessageHandler{Message: constants.UpdateProfile, ErrCode: model.NoError})
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
