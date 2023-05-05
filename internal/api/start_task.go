package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type StartTaskRequest struct {
	TeamRequest
	TaskRequest
}

func (r *StartTaskRequest) Validate(requestCtx context.Context) []error {
	errs := r.TeamRequest.Validate(requestCtx)
	return append(errs, r.TaskRequest.Validate(requestCtx)...)
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
