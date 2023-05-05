package api

import (
	"context"
	"net/http"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

type ShowTaskHintRequest struct {
	TeamRequest
	TaskRequest
	HintNum int `json:"hint_num"`
}

func (r *ShowTaskHintRequest) Validate(requestCtx context.Context) []error {
	errs := r.TeamRequest.Validate(requestCtx)
	errs = append(errs, r.TaskRequest.Validate(requestCtx)...)
	if r.HintNum < 0 {
		errs = append(errs, utils.NewErrWithType(errors.New("hint_num"), ErrorInvalidRequest))
	}

	return errs
}

type ShowTaskHintResponse struct {
	NextHintNum int    `json:"next_num"`
	Hint        string `json:"hint"`
}

func (c *Controller) ShowTaskHint(w http.ResponseWriter, r *http.Request) {
	req := &ShowTaskHintRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	nextHintNum, hint, err := c.service.ShowTaskHint(req.TeamId, req.TaskId, req.HintNum)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "ShowTaskHint"))
		return
	}
	writeSuccessResponse(w, ShowTaskHintResponse{NextHintNum: nextHintNum, Hint: hint})
}
