package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
)

func (s PgRepo) StartTask(teamId int, taskId int) error {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return err
	}
	if s.checkContestExist(contest) != nil {
		return err
	}
	if s.checkContestStarting(contest) != nil {
		return err
	}
	if s.checkContestFinished(contest) != nil {
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
	if err != nil {
		return err
	}
	if s.checkTaskExist(task) != nil {
		return err
	}
	if s.checkTaskAlreadyStarted(teamTask) != nil {
		return err
	}

	query := `
		insert into team_task 
		    (team_id, task_id, answers, answer_uuids, next_hint_num, status)
		values ($1, $2, '{}', '{}', 0, $3)
	`
	res, err := s.db.Exec(query, teamId, taskId, api.TaskStatusStarted)
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
