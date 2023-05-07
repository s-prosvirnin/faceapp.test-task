package api

import (
	"context"
	"net/http"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

type SendTaskAnswerRequest struct {
	TeamRequest
	TaskRequest
	Answer     string `json:"answer"`
	AnswerUuid string `json:"answer_uuid"`
}

func (r *SendTaskAnswerRequest) Validate(requestCtx context.Context) []error {
	errs := r.TeamRequest.Validate(requestCtx)
	errs = append(errs, r.TaskRequest.Validate(requestCtx)...)
	if r.Answer == "" {
		errs = append(errs, utils.NewErrWithType(errors.New("answer"), ErrorTypeInvalidRequest))
	}
	if r.AnswerUuid == "" {
		errs = append(errs, utils.NewErrWithType(errors.New("answer_uuid"), ErrorTypeInvalidRequest))
	}

	return errs
}

type SendTaskAnswerResponse struct {
	AnswerPassed bool `json:"answer_passed"`
}

func (c *Controller) SendTaskAnswer(w http.ResponseWriter, r *http.Request) {
	req := &SendTaskAnswerRequest{}
	if !createRequestModelWithValidate(req, w, r) {
		return
	}

	answerPassed, err := c.service.SendTaskAnswer(req.TeamId, req.TaskId, req.Answer, req.AnswerUuid)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "SendTaskAnswer"))
		return
	}
	writeSuccessResponse(w, SendTaskAnswerResponse{AnswerPassed: answerPassed})
}
