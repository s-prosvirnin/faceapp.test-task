package api

import (
	"net/http"

	"github.com/pkg/errors"
)

type StartTaskRequest struct {
	TeamRequest
	TaskRequest
}

func (r *StartTaskRequest) Validate() []error {
	errs := r.TeamRequest.Validate()
	return append(errs, r.TeamRequest.Validate()...)
}

func (c *Controller) StartTask(w http.ResponseWriter, r *http.Request) {
	req := &StartTaskRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	err := c.service.StartTask(req.TeamId, req.TaskId)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "StartTask"))
		return
	}
	writeSuccessResponse(w, true)
}
