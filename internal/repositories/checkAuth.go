package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
)

func (r PgRepo) GetTeamIdByAuthToken(accessToken string) (int, error) {
	query := `
		select a.id
		from teamEntity a
		where access_token = $1
	`

	var teamId int
	err := r.db.Get(&teamId, query, accessToken)
	if err == sql.ErrNoRows {
		return 0, utils.NewErrWithType(api.ErrAuthTokenInvalid, api.ErrorTypeDomain)
	}
	if err != nil {
		return 0, wrapInternalError(err, "db.Get")
	}

	return teamId, nil
}
