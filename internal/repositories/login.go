package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/google/uuid"
)

func (s PgRepo) Login(login string, password string) (api.GetAuthResponse, error) {
	var team team
	// @todo: в теории тут может быть паника
	authToken := uuid.New().String()

	query := `
		update team set access_token = $1
		where login = $2 and password = $3
		returning *
	`

	err := s.db.QueryRowx(query, authToken, login, password).StructScan(&team)
	if err == sql.ErrNoRows {
		return api.GetAuthResponse{}, utils.NewErrWithType(api.ErrLoginPasswordInvalid, api.ErrorDomainType)
	}
	if err != nil {
		return api.GetAuthResponse{}, wrapInternalError(err, "db.Get")
	}

	return api.GetAuthResponse{
		TeamId:    team.Id,
		AuthToken: team.AccessToken,
	}, nil
}
