package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

const ErrorInternalType = "internal_error"
const ErrorDomainType = "domain_error"
const ErrorInvalidRequest = "invalid_request"

var ErrAuthTokenInvalid = errors.New("auth_token_invalid")
var ErrTeamNotFound = errors.New("team_not_found")
var ErrContestNotFound = errors.New("contest_not_found")
var ErrContestNotStarted = errors.New("contest_not_started")
var ErrContestFinished = errors.New("contest_finished")
var ErrTaskNotFound = errors.New("task_not_found")
var ErrTaskAlreadyStarted = errors.New("task_already_started")
var ErrTaskNotStarted = errors.New("task_not_started")
var ErrTaskHintNumNotExist = errors.New("hint_num_not_exist")

// @todo: pros возможно, нам не нужны статусы в БД, т.к. можно создавать запись в момент начала задания
const TaskStatusPassed = "passed"
const TaskStatusNotStarted = "not_started"
const TaskStatusStarted = "started"
const TaskStatusAttemptFailed = "attempt_failed"

const ContestStatusActual = "actual"

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

type Service interface {
	GetContest(teamId int) (GetContestResponse, error)
	Login(login string, password string) (GetAuthResponse, error)
	CheckAuth(accessToken string) error
	GetTeamTasks(teamId int) (GetTasksResponse, error)
	StartTask(teamId int, taskId int) error
	SendTaskAnswer(teamId int, taskId int, teamAnswer string, answerUuid string) (
		answerPassed bool,
		err error,
	)
	ShowTaskHint(teamId int, taskId int, hintNum int) (nextHintNum int, hint string, err error)
	GetContestResults(teamId int) (GetContestResultsResponse, error)
}

type Validatable interface {
	Validate() []error
}

func validateRequest(reqModel Validatable, writer http.ResponseWriter, request *http.Request) bool {
	if err := json.NewDecoder(request.Body).Decode(reqModel); err != nil {
		writeErrorResponse(
			writer,
			utils.NewErrWithType(errors.Wrap(err, "validateRequest json.NewDecoder"), ErrorInternalType),
		)

		return false
	}
	if errs := reqModel.Validate(); len(errs) > 0 {
		writeValidationErrorResponse(writer, errs)

		return false
	}

	return true
}

type errorsResponse struct {
	Errors errorResponse `json:"errors"`
}

type errorResponse struct {
	Keys                 []string `json:"keys"`
	InvalidRequestFields []string `json:"invalid_request_fields"`
}

type successResponse struct {
	Success interface{} `json:"success"`
}

func newSuccessResponse(jsonableResponse any) successResponse {
	return successResponse{
		Success: jsonableResponse,
	}
}

func writeSuccessResponse(writer http.ResponseWriter, jsonableResponse any) {
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(newSuccessResponse(jsonableResponse))
	if err != nil {
		log.Print(errors.Wrap(err, "json.NewEncoder.newSuccessResponse"))
	}
}

func writeErrorResponse(writer http.ResponseWriter, err error) {
	// разделяем вывод ошибок на внутренние и внешние
	errorType := utils.GetErrorType(err)
	switch errorType.Type() {
	case ErrorInternalType:
		log.Print(err)
		err = errors.New(ErrorInternalType)
	case ErrorDomainType:
		err = errors.New(errorType.Error())
	}

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(
		errorsResponse{
			Errors: errorResponse{Keys: []string{err.Error()}, InvalidRequestFields: []string{}},
		},
	)
	if err != nil {
		log.Print(errors.Wrap(err, "writeErrorResponse.json.NewEncoder"))
	}
}

func writeValidationErrorResponse(writer http.ResponseWriter, errs []error) {
	writer.WriteHeader(http.StatusOK)

	errorAsStrings := make([]string, 0, len(errs))
	for _, err := range errs {
		errorAsStrings = append(errorAsStrings, err.Error())
	}

	response := errorsResponse{
		Errors: errorResponse{Keys: []string{ErrorInvalidRequest}, InvalidRequestFields: errorAsStrings},
	}

	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Print(errors.Wrap(err, "writeValidationErrorResponse.json.NewEncoder"))
	}
}

type TeamRequest struct {
	TeamId int `json:"team_id" validate:"required" minimum:"1"`
}

func (r *TeamRequest) Validate() []error {
	var errs []error
	if r.TeamId <= 0 {
		errs = append(errs, errors.New("team_id"))
	}

	return errs
}

type TaskRequest struct {
	TaskId int `json:"task_id" validate:"required" minimum:"1"`
}

func (r *TaskRequest) Validate() []error {
	var errs []error
	if r.TaskId <= 0 {
		errs = append(errs, errors.New("task_id"))
	}

	return errs
}
