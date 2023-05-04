package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
)

func (s PgRepo) CheckAuth(accessToken string) error {
	query := `
		select a.id
		from team a
		where access_token = $1
	`

	var teamId int
	err := s.db.Get(&teamId, query, accessToken)
	if err == sql.ErrNoRows {
		return utils.NewErrWithType(api.ErrAuthTokenInvalid, api.ErrorDomainType)
	}
	if err != nil {
		return wrapInternalError(err, "db.Get")
	}

	return nil
}
