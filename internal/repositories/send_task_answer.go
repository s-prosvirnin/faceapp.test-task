package repositories

import (
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/lib/pq"
)

func (s PgRepo) SendTaskAnswer(teamId int, taskId int, teamAnswer string, answerUuid string) (
	answerPassed bool,
	err error,
) {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return false, err
	}
	if err = s.checkContestExist(contest); err != nil {
		return false, err
	}
	if err = s.checkContestStarting(contest); err != nil {
		return false, err
	}
	if err = s.checkContestFinished(contest); err != nil {
		return false, err
	}

	task, err := s.getTaskEntity(contest.Id, taskId)
	if err != nil {
		return false, err
	}
	teamTask, err := s.getTeamTaskEntity(teamId, taskId)
	if err != nil {
		return false, err
	}
	if err = s.checkTaskExist(task); err != nil {
		return false, err
	}
	if err = s.checkTaskNotStarted(teamTask); err != nil {
		return false, err
	}
	if isTaskPassed(task, teamTask) {
		return false, err
	}
	// @todo: pros добавить ограничения на кол-во

	// проверяем идемпотентность
	for _, uuid := range teamTask.AnswersUuid {
		if answerUuid == uuid {
			return isAnswerPassed(task, teamAnswer), nil
		}
	}

	teamTask.Answers = append(teamTask.Answers, teamAnswer)
	teamTask.AnswersUuid = append(teamTask.AnswersUuid, answerUuid)
	teamTask.AnswersCreatedAt = append(teamTask.AnswersCreatedAt, string(pq.FormatTimestamp(time.Now())))
	teamTask.Status = api.TaskStatusAttemptFailed
	if isAnswerPassed(task, teamAnswer) {
		teamTask.Status = api.TaskStatusPassed
	}

	query := `
		update team_task set 
		    answers = :answers
		    , answer_uuids = :answer_uuids
		    , answers_created_at = :answers_created_at
		    , status = :status
		where team_id = :team_id and task_id = :task_id
	`
	res, err := s.db.NamedExec(query, teamTask)
	if err != nil {
		return false, wrapInternalError(err, "db.Exec")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, wrapInternalError(err, "rows affected error")
	}
	if rowsAffected != 1 {
		return false, wrapInternalError(err, "rows affected  mismatch")
	}

	return isAnswerPassed(task, teamAnswer), nil
}