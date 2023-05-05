package api

import (
	"context"
	"net/http"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

func (a *AuthRequest) Validate(requestCtx context.Context) []error {
	var errs []error
	if a.Login == "" {
		errs = append(errs, utils.NewErrWithType(errors.New("login"), ErrorInvalidRequest))
	}
	if a.Password == "" {
		errs = append(errs, utils.NewErrWithType(errors.New("pass"), ErrorInvalidRequest))
	}

	return errs
}

type GetAuthResponse struct {
	AuthToken string `json:"auth_token"`
	TeamId    int    `json:"team_id"`
}

func (c *Controller) Auth(w http.ResponseWriter, r *http.Request) {
	req := &AuthRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	// упростил - пароль не шифруется, передается в открытом виде
	// упростил - используется только один бесконечный токен
	response, err := c.service.Login(req.Login, req.Password)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "AuthRequest"))
		return
	}

	writeSuccessResponse(w, response)
}
