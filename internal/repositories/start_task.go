package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
)

func (r PgRepo) StartTask(teamId int, taskId int) error {
	contest, err := r.getContestEntity(teamId)
	if err != nil {
		return err
	}
	if err = checkContestExist(contest); err != nil {
		return err
	}
	if err = checkContestStarting(contest); err != nil {
		return err
	}
	if err = checkContestFinished(contest); err != nil {
		return err
	}

	task, err := r.getTaskEntity(contest.Id, taskId)
	if err != nil {
		return err
	}
	teamTask, err := r.getTeamTaskEntity(teamId, taskId)
	if err != nil {
		return err
	}
	if err = checkTaskExist(task); err != nil {
		return err
	}
	if err = checkTaskAlreadyStarted(teamTask); err != nil {
		return err
	}

	teamTask.TeamId = teamId
	teamTask.TaskId = taskId
	teamTask.Answers = []string{}
	teamTask.AnswersUuid = []string{}
	teamTask.AnswersCreatedAt = []string{}
	teamTask.NextHintNum = 0
	teamTask.Status = api.TaskStatusStarted

	query := `
		insert into team_task 
		    (team_id, task_id, answers, answer_uuids, answers_created_at, next_hint_num, status)
		values (:team_id, :task_id, :answers, :answer_uuids, :answers_created_at, :next_hint_num, :status)
	`
	res, err := r.db.NamedExec(query, teamTask)
	if err != nil {
		return wrapInternalError(err, "db.NamedExec")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return wrapInternalError(err, "rows affected error")
	}
	if rowsAffected != 1 {
		return wrapInternalError(err, "rows affected  mismatch")
	}

	return nil
}
