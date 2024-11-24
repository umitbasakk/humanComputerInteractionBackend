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
  Users (name,username, password,email,phone,token)
VALUES
  ($1, $2,$3,$4,$5,$6)
`

const createVerify = `-- name: CreateVerify :exec
INSERT INTO
  VerifyUsers (username,verify_code,verify_status)
VALUES
  ($1, $2,$3)
`
const getVerifyCodeByUsername = `-- name: GetVerifyCodeByUsername :one
SELECT
 *
FROM
  VerifyUsers
WHERE
  username = $1
LIMIT
  1
`

const GetUserByUsername = `-- name: GetUserByUsername :one
SELECT
 *
FROM
  users
WHERE
  username = $1
LIMIT
  1
`

const GetUserByEmail = `-- name: GetUserByEmail :one
SELECT
 *
FROM
  users
WHERE
  email = $1
LIMIT
  1
`
const updateVerify = `-- name: UpdateVerify :exec
UPDATE
  VerifyUsers
SET
  verify_status = $1
WHERE
  username = $2
`

const updateVerifyCode = `-- name: UpdateVerifyCode :exec
UPDATE
  VerifyUsers
SET
  verify_code = $1
WHERE
  username = $2
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
	_, err := dl.connPs.Query(createUser, user.Name, user.Username, user.Password, user.Email, user.Phone, user.Token)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) Login(ctx echo.Context, username string) (*model.User, error) {
	user := &model.User{}
	result, err := dl.connPs.Query(GetUserByUsername, username)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
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
	result, err := dl.connPs.Query(GetUserByUsername, username)
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

func (dl *UserDatalayerImpl) GetVerifyCode(ctx echo.Context, username string) (*model.Verify, error) {
	vf := &model.Verify{}
	result, err := dl.connPs.Query(getVerifyCodeByUsername, username)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&vf.Id, &vf.Username, &vf.VerifyCode, &vf.VerifyStatus, &vf.Created_at, &vf.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}
	return vf, nil
}

func (dl *UserDatalayerImpl) CreateVerifyCode(ctx echo.Context, verify *model.Verify) error {
	_, err := dl.connPs.Query(createVerify, verify.Username, verify.VerifyCode, 0)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) VerifyCode(ctx echo.Context, username string) error {
	_, err := dl.connPs.Query(updateVerify, 1, username)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) GetUserByEmail(ctx echo.Context, email string) (*model.User, error) {
	user := &model.User{}
	result, err := dl.connPs.Query(GetUserByEmail, email)
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

func (dl *UserDatalayerImpl) UpdateVerifyCode(ctx echo.Context, username string, vCode string) error {
	_, err := dl.connPs.Query(updateVerifyCode, vCode, username)
	if err != nil {
		return err
	}
	return nil
}
