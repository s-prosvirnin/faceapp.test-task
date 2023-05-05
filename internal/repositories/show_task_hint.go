package repositories

func (s PgRepo) ShowTaskHint(teamId int, taskId int, hintNum int) (nextHintNum int, hint string, err error) {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return 0, "", err
	}
	if err = s.checkContestExist(contest); err != nil {
		return 0, "", err
	}
	if err = s.checkContestStarting(contest); err != nil {
		return 0, "", err
	}
	if err = s.checkContestFinished(contest); err != nil {
		return 0, "", err
	}

	task, err := s.getTaskEntity(contest.Id, taskId)
	if err != nil {
		return 0, "", err
	}
	teamTask, err := s.getTeamTaskEntity(teamId, taskId)
	if err != nil {
		return 0, "", err
	}
	if err = s.checkTaskExist(task); err != nil {
		return 0, "", err
	}
	if err = s.checkTaskNotStarted(teamTask); err != nil {
		return 0, "", err
	}
	if err = s.checkTaskNextHintNumExist(task, hintNum); err != nil {
		return 0, "", err
	}
	// если номер подсказки для показа меньше чем текущий, то показываем последнюю показанную подсказку
	if s.isTaskHintAlreadyShown(teamTask, hintNum) {
		return getHintNumForResponse(task, teamTask), getHint(task, teamTask.NextHintNum), nil
	}

	teamTask.NextHintNum++

	query := `
		update team_task
		set next_hint_num = :next_hint_num
		where team_id = :team_id and task_id = :task_id
	`
	res, err := s.db.NamedExec(query, teamTask)
	if err != nil {
		return 0, "", wrapInternalError(err, "db.Exec")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, "", wrapInternalError(err, "rows affected error")
	}
	if rowsAffected != 1 {
		return 0, "", wrapInternalError(err, "rows affected  mismatch")
	}

	return getHintNumForResponse(task, teamTask), getHint(task, teamTask.NextHintNum), nil
}

func getHint(task taskEntity, nextNum int) string {
	if isNextHintNumLast(task, nextNum) {
		return task.Hints[len(task.Hints)-1]
	}

	return task.Hints[nextNum-1]
}
