package database

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO
  users (name,username,email, phone,password,token)
VALUES
  ($1,$2,$3,$4,$5,$6)
`

const createVerify = `-- name: CreateVerify :exec
INSERT INTO
  VerifyUsers (user_id,verify_code,verify_status)
VALUES
  ($1, $2,$3)
`
const getVerifyCodeByUserId = `-- name: getVerifyCodeByUserId :one
SELECT
 *
FROM
  verifyusers
WHERE
  user_id = $1
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

const GetUserByPhone = `-- name: GetUserByUsername :one
SELECT
 *
FROM
  users
WHERE
  phone = $1
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
  user_id = $2
`

const updateVerifyCode = `-- name: UpdateVerifyCode :exec
UPDATE
  VerifyUsers
SET
  verify_code = $1,
  updated_at = NOW()
WHERE
  user_id = $2
`

const updateProfile = `-- name: UpdateProfile :exec
UPDATE
  users
SET
  name = $1,
  username = $2,
  email = $3
WHERE
  username = $4
`

const updatePassword = `-- name: UpdatePassword :exec
UPDATE
  users
SET
  password = $1
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

func (dl *UserDatalayerImpl) Signup(tx *sql.Tx, ctx echo.Context, user *model.User) error {
	rows, err := tx.Query(createUser, user.Name, user.Username, user.Email, user.Phone, user.Password, user.Token)
	if err != nil {
		return err
	}
	defer rows.Close()
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
	defer result.Close()
	return user, nil
}

func (dl *UserDatalayerImpl) GetUserUsername(tx *sql.Tx, ctx echo.Context, username string) (*model.User, error) {
	user := &model.User{}
	result, err := tx.Query(GetUserByUsername, username)
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
func (dl *UserDatalayerImpl) GetUserEmail(tx *sql.Tx, ctx echo.Context, email string) (*model.User, error) {
	user := &model.User{}
	result, err := tx.Query(GetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}

	}
	defer result.Close()
	return user, nil
}

func (dl *UserDatalayerImpl) SaveTokenByUsername(tx *sql.Tx, ctx echo.Context, username string, token string) error {
	rows, err := tx.Query(updateToken, username, token)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (dl *UserDatalayerImpl) IsThereEqualUsername(tx *sql.Tx, ctx echo.Context, username string) error {
	result, err := tx.Query(GetUserByUsername, username)
	if err != nil {
		return err
	}
	if result.Next() {
		return errors.New("Already used email")
	}
	defer result.Close()
	return nil
}

func (dl *UserDatalayerImpl) GetUserByPhone(tx *sql.Tx, ctx echo.Context, phone string) error {
	result, err := tx.Query(GetUserByPhone, phone)
	if err != nil {
		return err
	}
	if result.Next() {
		return errors.New("Already used phone")
	}
	result.Close()
	return nil

}

func (dl *UserDatalayerImpl) GetVerifyCode(tx *sql.Tx, ctx echo.Context, user_id int) (*model.Verify, error) {
	vf := &model.Verify{}
	result, err := tx.Query(getVerifyCodeByUserId, strconv.Itoa(user_id))
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&vf.Id, &vf.Username, &vf.VerifyCode, &vf.VerifyStatus, &vf.Created_at, &vf.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}
	defer result.Close()
	return vf, nil
}

func (dl *UserDatalayerImpl) CreateVerifyCode(tx *sql.Tx, ctx echo.Context, verify *model.Verify, user_id int) error {
	_, err := tx.Query(createVerify, strconv.Itoa(user_id), verify.VerifyCode, 0)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) VerifyCode(tx *sql.Tx, ctx echo.Context, user_id int) error {
	rows, err := tx.Query(updateVerify, 1, strconv.Itoa(user_id))
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (dl *UserDatalayerImpl) IsThereEqualEmail(tx *sql.Tx, ctx echo.Context, email string) error {
	result, err := tx.Query(GetUserByEmail, email)
	if err != nil {
		return err
	}
	if result.Next() {
		return errors.New("Already used email")
	}
	result.Close()
	return nil
}

func (dl *UserDatalayerImpl) UpdateVerifyCode(tx *sql.Tx, ctx echo.Context, user_id int, vCode string) error {
	rows, err := tx.Query(updateVerifyCode, vCode, strconv.Itoa(user_id))
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (dl *UserDatalayerImpl) ChangePassword(ctx echo.Context, username string, password string) error {
	_, err := dl.connPs.Query(updatePassword, password, username)
	if err != nil {
		return err
	}
	return nil
}

func (dl *UserDatalayerImpl) UpdateProfile(tx *sql.Tx, ctx echo.Context, profile *model.UpdateProfileRequest, username string) error {
	rows, err := tx.Query(updateProfile, profile.Name, profile.Username, profile.Email, username)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (dl *UserDatalayerImpl) GetTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := dl.connPs.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (dl *UserDatalayerImpl) CommitTransaction(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
