package middlewares

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/database"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
	"github.com/umitbasakk/humanComputerInteractionBackend/helpers"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *sql.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") == false {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		accessToken := strings.Split(authHeader, " ")[1]

		claims, err := helpers.ParseJWT(accessToken)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		if helpers.IsClaimExpired(claims) {
			return c.JSON(http.StatusUnauthorized, "Expired")
		}
		counter := 0
		user := &model.User{}
		result, err := appMiddleware.DB.Query(database.GetUserByUsername, claims.Username)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}

		if result.Next() {
			counter++
			errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
			if errLogin != nil {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}
		}
		if counter == 0 {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		c.Set("user", user)
		return next(c)
	}
}
