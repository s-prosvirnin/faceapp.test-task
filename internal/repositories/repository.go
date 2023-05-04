package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	teamTaskStatusActual = "actual"
)

type PgRepo struct {
	db *sqlx.DB
}

func NewPgRepo(db *sqlx.DB) PgRepo {
	return PgRepo{db: db}
}

func wrapInternalError(err error, wrapMessage string) error {
	return utils.NewErrWithType(errors.Wrap(err, wrapMessage), api.ErrorInternalType)
}
