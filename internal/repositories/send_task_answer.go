package repositories

import (
	"time"

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
	if s.checkContestExist(contest) != nil {
		return false, err
	}
	if s.checkContestStarting(contest) != nil {
		return false, err
	}
	if s.checkContestFinished(contest) != nil {
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
	if s.checkTaskExist(task) != nil {
		return false, err
	}
	if s.checkTaskNotStarted(teamTask) != nil {
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
	teamTask.AnswersUuid = append(teamTask.Answers, answerUuid)
	teamTask.AnswersCreatedAt = append(teamTask.Answers, string(pq.FormatTimestamp(time.Now())))

	query := `
		update team_task set 
		    answers = :answers
		    , answer_uuids = :answer_uuids
		    , answers_created_at = :answers_created_at
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
