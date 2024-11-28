package database

import (
	"database/sql"

	"github.com/umitbasakk/humanComputerInteractionBackend/interfaces"
)

type AIDataLayerImp struct {
	connPs *sql.DB
}

func NewAIDataLayerImpl(db *sql.DB) interfaces.AIDataLayer {
	return &AIDataLayerImp{
		connPs: db,
	}
}
