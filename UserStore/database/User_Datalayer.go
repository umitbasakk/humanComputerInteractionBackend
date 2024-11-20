package database

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO
  Users (name,username, password,email,token)
VALUES
  ($1, $2,$3,$4,$5)
`

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT
 *
FROM
  users
WHERE
  username = $1
LIMIT
  1
`
const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
 *
FROM
  users
WHERE
  email = $1
LIMIT
  1
`

const updateToken = `-- name: UpdateToken :exec
UPDATE
  users
SET
  token = $1
WHERE
  username = $2
`

type UserDatalayerImpl struct {
	userDL interfaces.UserDatalayer
	connPs *sql.DB
}

func NewUserDatalayerImpl(conn *sql.DB) interfaces.UserDatalayer {
	return &UserDatalayerImpl{
		connPs: conn,
	}
}

func (dl *UserDatalayerImpl) GetUserByID(ctx context.Context, userID int16) *model.User {
	return nil
}

func (dl *UserDatalayerImpl) Signup(ctx echo.Context, user *model.User) error {
	_, err := dl.connPs.Query(createUser, user.Name, user.Username, user.Password, user.Email, user.Token)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) Login(ctx echo.Context, username string) (*model.User, error) {
	user := &model.User{}
	result, err := dl.connPs.Query(getUserByUsername, username)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}

	return user, nil
}

func (dl *UserDatalayerImpl) SaveTokenByUsername(ctx echo.Context, username string, token string) error {
	_, err := dl.connPs.Query(updateToken, username, token)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) GetUserByUsername(ctx echo.Context, username string) (*model.User, error) {
	user := &model.User{}
	result, err := dl.connPs.Query(getUserByUsername, username)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}
	return user, nil

}

func (dl *UserDatalayerImpl) GetUserByEmail(ctx echo.Context, email string) (*model.User, error) {
	user := &model.User{}
	result, err := dl.connPs.Query(getUserByEmail, email)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}
	return user, nil
}
