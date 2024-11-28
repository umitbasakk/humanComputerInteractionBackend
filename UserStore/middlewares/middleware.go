package middlewares

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/database"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	"github.com/umitbasakk/humanComputerInteractionBackend/constants"
	"github.com/umitbasakk/humanComputerInteractionBackend/helpers"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *sql.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") == false {
			return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
		}
		accessToken := strings.Split(authHeader, " ")[1]

		claims, err := helpers.ParseJWT(accessToken)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
		}
		if helpers.IsClaimExpired(claims) {
			return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
		}
		counter := 0
		user := &model.User{}
		result, err := appMiddleware.DB.Query(database.GetUserByUsername, claims.Username)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
		}

		if result.Next() {
			counter++
			errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
			if errLogin != nil {
				return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
			}
		}
		if counter == 0 {
			return c.JSON(http.StatusUnauthorized, &model.MessageHandler{Message: constants.UnauthorizedRequest, ErrCode: model.Authorized})
		}
		c.Set("user", user)
		return next(c)
	}
}
