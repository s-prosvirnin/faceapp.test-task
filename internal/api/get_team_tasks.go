package api

import (
	"net/http"

	"github.com/pkg/errors"
)

type GetTasksRequest struct {
	TeamRequest
}

type GetTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Coords      CoordsResponse   `json:"coords"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Answers     []AnswerResponse `json:"answers"`
	Hints       HintsResponse    `json:"hints"`
}

type AnswerResponse struct {
	Answer   string `json:"answer"`
	IsPassed bool   `json:"is_passed"`
}

type HintsResponse struct {
	Opened  []string `json:"opened"`
	Total   int      `json:"total"`
	NextNum int      `json:"next_num"`
}

type CoordsResponse struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (c *Controller) GetTeamTasks(w http.ResponseWriter, r *http.Request) {
	req := &GetContestRequest{}
	if !createRequestModelWithValidate(req, w, r) {
		return
	}

	response, err := c.service.GetTeamTasks(req.TeamId)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "GetTeamTasks"))
		return
	}
	writeSuccessResponse(w, response)
}
