package database

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	model "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/model/AI"
	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

const createRequest = `-- name: CreateRequest :exec
INSERT INTO
  requests (user_id,started_date,end_date, hash_tag,category,quantity_limit,request_status)
VALUES
  ($1,$2,$3,$4,$5,$6,$7)
`

const getRequestsOfUser = `-- name: GetRequestsOfUser :many
SELECT
 *
FROM
  requests
WHERE
  user_id = $1
ORDER BY created_at DESC
`

type AIDataLayerImp struct {
	connPs *sql.DB
}

func NewAIDataLayerImpl(db *sql.DB) interfaces.AIDataLayer {
	return &AIDataLayerImp{
		connPs: db,
	}
}

func (dl *AIDataLayerImp) GetTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := dl.connPs.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (dl *AIDataLayerImp) CommitTransaction(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (dl *AIDataLayerImp) SaveAiRequest(tx *sql.Tx, ctx echo.Context, aiData *model.AIData) error {
	row, err := tx.Query(createRequest, aiData.UserId, aiData.StartedDate, aiData.EndDate, aiData.HashTag, aiData.Category, aiData.QuantityLimit, aiData.RequestStatus)
	if err != nil {
		return err
	}
	defer row.Close()
	return nil
}

func (dl *AIDataLayerImp) GetRequestOfUser(tx *sql.Tx, ctx echo.Context, user_id string) ([]model.AIData, error) {
	requests := make([]model.AIData, 0)
	requestTmp := &model.AIData{}
	unUsedId := 0
	rows, err := tx.Query(getRequestsOfUser, user_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&unUsedId, &requestTmp.UserId, &requestTmp.StartedDate, &requestTmp.EndDate, &requestTmp.HashTag, &requestTmp.Category, &requestTmp.QuantityLimit, &requestTmp.RequestStatus, &requestTmp.Created_at, &requestTmp.Updated_at)
		if err != nil {
			return nil, err
		}
		requests = append(requests, *requestTmp)
	}
	defer rows.Close()
	return requests, nil
}
