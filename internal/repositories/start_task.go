package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
)

func (s PgRepo) StartTask(teamId int, taskId int) error {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return err
	}
	if err = s.checkContestExist(contest); err != nil {
		return err
	}
	if err = s.checkContestStarting(contest); err != nil {
		return err
	}
	if err = s.checkContestFinished(contest); err != nil {
		return err
	}

	task, err := s.getTaskEntity(contest.Id, taskId)
	if err != nil {
		return err
	}
	teamTask, err := s.getTeamTaskEntity(teamId, taskId)
	if err != nil {
		return err
	}
	if err = s.checkTaskExist(task); err != nil {
		return err
	}
	if err = s.checkTaskAlreadyStarted(teamTask); err != nil {
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
	res, err := s.db.NamedExec(query, teamTask)
	if err != nil {
		return wrapInternalError(err, "db.Exec")
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
