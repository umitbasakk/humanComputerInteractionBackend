package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/Auth"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO
  users (name,username,email,password,token)
VALUES
  ($1,$2,$3,$4,$5)
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
	rows, err := tx.Query(createUser, user.Name, user.Username, user.Email, user.Password, user.Token)
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
	defer result.Close()

	if !result.Next() {
		return nil, fmt.Errorf("user not found with username: %s", username)
	}

	err = result.Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		return nil, fmt.Errorf("error scanning user data: %v", err)
	}

	return user, nil
}

func (dl *UserDatalayerImpl) GetUserUsername(tx *sql.Tx, ctx echo.Context, username string) (*model.User, error) {
	user := &model.User{}
	result, err := tx.Query(GetUserByUsername, username)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
		if errLogin != nil {
			return nil, errLogin
		}
	}
	defer result.Close()
	return user, nil
}
func (dl *UserDatalayerImpl) GetUserEmail(tx *sql.Tx, ctx echo.Context, email string) (*model.User, error) {
	user := &model.User{}
	result, err := tx.Query(GetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	if result.Next() {
		errLogin := result.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Token, &user.Created_at, &user.Updated_at)
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
	defer rows.Close()
	return nil
}

func (dl *UserDatalayerImpl) IsThereEqualUsername(tx *sql.Tx, ctx echo.Context, username string) error {
	result, err := tx.Query(GetUserByUsername, username)
	if err != nil {
		return err
	}
	defer result.Close()
	if result.Next() {
		return errors.New("Already used email")
	}
	return nil
}

func (dl *UserDatalayerImpl) IsThereEqualEmail(tx *sql.Tx, ctx echo.Context, email string) error {
	result, err := tx.Query(GetUserByEmail, email)
	if err != nil {
		return err
	}
	defer result.Close()
	if result.Next() {
		return errors.New("Already usedd email")
	}
	return nil
}

func (dl *UserDatalayerImpl) ChangePassword(ctx echo.Context, username string, password string) error {
	Ps, err := dl.connPs.Query(updatePassword, password, username)
	if err != nil {
		return err
	}
	defer Ps.Close()
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
