package api

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type GetContestRequest struct {
	TeamRequest
}

type GetContestResponse struct {
	ContestId int       `json:"contest_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	// @todo: pros добавить поля
}

func (c *Controller) GetContest(w http.ResponseWriter, r *http.Request) {
	req := &GetContestRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	response, err := c.service.GetContest(req.TeamId)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "GetContest"))
		return
	}
	// @todo: pros добавить поля в response
	writeSuccessResponse(w, response)
}
