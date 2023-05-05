package api

import (
	"net/http"

	"github.com/pkg/errors"
)

type GetContestResultsRequest struct {
	TeamRequest
}

type GetContestResultsResponse struct {
	TeamResults []TeamResultResponse `json:"teams_results"`
}

type TeamResultResponse struct {
	TeamRank         int                  `json:"team_rank"`
	TeamName         string               `json:"team_name"`
	TaskResults      []TaskResultResponse `json:"tasks_results"`
	TasksPassedCount int                  `json:"tasks_passed_count"`
	PenaltyTimeSec   int                  `json:"penalty_time_sec"`
}

type TaskResultResponse struct {
	TaskId           int    `json:"task_id"`
	Status           string `json:"status"`
	HintsOpenedCount int    `json:"hints_opened_count"`
}

func (c *Controller) GetContestResults(w http.ResponseWriter, r *http.Request) {
	req := &GetContestResultsRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	response, err := c.service.GetContestResults(req.TeamId)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "GetContestResults"))
		return
	}
	writeSuccessResponse(w, response)
}
