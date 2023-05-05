package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

const ErrorInternalType = "internal_error"
const ErrorDomainType = "domain_error"
const ErrorInvalidRequest = "invalid_request"

const TaskStatusPassed = "passed"
const TaskStatusNotStarted = "not_started"
const TaskStatusStarted = "started"
const TaskStatusAttemptFailed = "attempt_failed"

const ContestStatusStarted = "started"
const ContestStatusCompleted = "completed"
const ContestStatusWillStartSoon = "will_start_soon"

const teamIdContextKey = "team_id"

var ErrLoginPasswordInvalid = errors.New("login_pass_invalid")
var ErrAuthTokenInvalid = errors.New("auth_token_invalid")

var ErrTeamNotFound = errors.New("team_not_found")
var ErrContestNotFound = errors.New("contest_not_found")
var ErrContestNotStarted = errors.New("contest_not_started")
var ErrContestFinished = errors.New("contest_finished")
var ErrTaskNotFound = errors.New("task_not_found")
var ErrTaskAlreadyStarted = errors.New("task_already_started")
var ErrTaskAnswerAlreadyPassed = errors.New("answer_already_passed")
var ErrTaskNotStarted = errors.New("task_not_started")
var ErrTaskHintNumNotExist = errors.New("hint_num_not_exist")
var ErrTaskAnswerPerTimeLimitExceeded = errors.New("answer_per_time_limit_exceeded")
var ErrTaskAnswersLimitExceeded = errors.New("answer_limit_exceeded")

type Controller struct {
	service Service
}

// @todo: pros разобраться с датами, go ставит текущую локаль

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

type Service interface {
	GetContest(teamId int) (GetContestResponse, error)
	Login(login string, password string) (GetAuthResponse, error)
	GetTeamIdByAuthToken(accessToken string) (int, error)
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
	Validate(requestCtx context.Context) []error
}

func validateRequest(reqModel Validatable, writer http.ResponseWriter, request *http.Request) bool {
	if err := json.NewDecoder(request.Body).Decode(reqModel); err != nil {
		writeErrorResponse(
			writer,
			utils.NewErrWithType(errors.Wrap(err, "validateRequest json.NewDecoder"), ErrorInternalType),
		)

		return false
	}
	if errs := reqModel.Validate(request.Context()); len(errs) > 0 {
		writeErrorsResponse(writer, errs)

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
	writeErrorsResponse(writer, []error{err})
}

func writeErrorsResponse(writer http.ResponseWriter, errs []error) {
	writer.WriteHeader(http.StatusOK)

	var validationErrorsKeys []string
	var errorsKeys []string
	for _, err := range errs {
		validationErrorsKeys = append(validationErrorsKeys, err.Error())

		// разделяем вывод ошибок на внутренние и внешние
		errorType := utils.GetErrorType(err)
		switch errorType.Type() {
		case ErrorInternalType:
			log.Print(err)
			errorsKeys = append(errorsKeys, ErrorInternalType)
		case ErrorDomainType:
			errorsKeys = append(errorsKeys, errorType.Error())
		case ErrorInvalidRequest:
			errorsKeys = append(errorsKeys, ErrorInvalidRequest)
			validationErrorsKeys = append(validationErrorsKeys, err.Error())
		}
	}

	response := errorsResponse{
		Errors: errorResponse{Keys: errorsKeys, InvalidRequestFields: validationErrorsKeys},
	}

	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Print(errors.Wrap(err, "writeErrorsResponse.json.NewEncoder"))
	}
}

type TeamRequest struct {
	TeamId int `json:"team_id"`
}

func (r *TeamRequest) Validate(requestCtx context.Context) []error {
	var errs []error
	if r.TeamId <= 0 {
		errs = append(errs, utils.NewErrWithType(errors.New("team_id"), ErrorInvalidRequest))
	}
	// @todo: убрать team_id из запросов, а брать его по токену
	ctxTeamId := requestCtx.Value(teamIdContextKey)
	if ctxTeamId != r.TeamId {
		errs = append(errs, utils.NewErrWithType(ErrTeamNotFound, ErrorDomainType))
	}

	return errs
}

type TaskRequest struct {
	TaskId int `json:"task_id"`
}

func (r *TaskRequest) Validate(requestCtx context.Context) []error {
	var errs []error
	if r.TaskId <= 0 {
		errs = append(errs, utils.NewErrWithType(errors.New("task_id"), ErrorInvalidRequest))
	}

	return errs
}
