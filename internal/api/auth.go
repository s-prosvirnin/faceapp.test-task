package api

import (
	"net/http"

	"github.com/pkg/errors"
)

var ErrLoginPasswordInvalid = errors.New("login_pass_invalid")

type AuthRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"pass" validate:"required"`
}

func (a *AuthRequest) Validate() []error {
	var errs []error
	if a.Login == "" {
		errs = append(errs, errors.New("login"))
	}
	if a.Password == "" {
		errs = append(errs, errors.New("password"))
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
	response, err := c.service.Login(req.Login, req.Password)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "AuthRequest"))
		return
	}

	writeSuccessResponse(w, response)
}
