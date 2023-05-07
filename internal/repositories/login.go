package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/google/uuid"
)

func (r PgRepo) Login(login string, password string) (api.GetAuthResponse, error) {
	team := teamEntity{
		Login:    login,
		Password: password,
	}

	query := `
		select * from team
		where login = $1 and password = $2
	`
	err := r.db.Get(&team, query, login, password)
	if err == sql.ErrNoRows {
		return api.GetAuthResponse{}, utils.NewErrWithType(api.ErrLoginPasswordInvalid, api.ErrorTypeDomain)
	}
	if err != nil {
		return api.GetAuthResponse{}, wrapInternalError(err, "db.Get")
	}

	if team.AccessToken.String != "" {
		return api.GetAuthResponse{
			TeamId:    team.Id,
			AuthToken: team.AccessToken.String,
		}, nil
	}

	// в теории тут может быть паника
	team.AccessToken.String = uuid.New().String()
	team.AccessToken.Valid = true
	query = `
		update team set access_token = :access_token
		where login = :login and password = :password
	`
	res, err := r.db.NamedExec(query, team)
	if rowsAffected, err := res.RowsAffected(); rowsAffected == 0 || err != nil {
		return api.GetAuthResponse{}, wrapInternalError(err, "rows affected mismatch")
	}
	if err != nil {
		return api.GetAuthResponse{}, wrapInternalError(err, "db.NamedExec")
	}

	return api.GetAuthResponse{
		TeamId:    team.Id,
		AuthToken: team.AccessToken.String,
	}, nil
}
