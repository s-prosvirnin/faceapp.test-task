package api

import (
	"net/http"

	"github.com/pkg/errors"
)

type ShowTaskHintRequest struct {
	TeamRequest
	TaskRequest
	HintNum int `json:"hint_num"`
}

func (r *ShowTaskHintRequest) Validate() []error {
	errs := r.TeamRequest.Validate()
	errs = append(errs, r.TeamRequest.Validate()...)
	if r.HintNum <= 0 {
		errs = append(errs, errors.New("hint_num"))
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
