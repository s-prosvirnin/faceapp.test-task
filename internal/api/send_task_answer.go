package api

import (
	"net/http"

	"github.com/pkg/errors"
)

type SendTaskAnswerRequest struct {
	TeamRequest
	TaskRequest
	Answer     string `json:"answer"`
	AnswerUuid string `json:"answer_uuid"`
}

func (r *SendTaskAnswerRequest) Validate() []error {
	errs := r.TeamRequest.Validate()
	errs = append(errs, r.TeamRequest.Validate()...)
	if r.Answer == "" {
		errs = append(errs, errors.New("answer"))
	}
	if r.AnswerUuid == "" {
		errs = append(errs, errors.New("answer_uuid"))
	}

	return errs
}

type SendTaskAnswerResponse struct {
	AnswerPassed bool `json:"answer_passed"`
}

func (c *Controller) SendTaskAnswer(w http.ResponseWriter, r *http.Request) {
	req := &SendTaskAnswerRequest{}
	if !validateRequest(req, w, r) {
		return
	}

	answerPassed, err := c.service.SendTaskAnswer(req.TeamId, req.TaskId, req.Answer, req.AnswerUuid)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "SendTaskAnswer"))
		return
	}
	writeSuccessResponse(w, SendTaskAnswerResponse{AnswerPassed: answerPassed})
}
