package repositories

func (s PgRepo) ShowTaskHint(teamId int, taskId int, hintNum int) (nextHintNum int, hint string, err error) {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return 0, "", err
	}
	if s.checkContestExist(contest) != nil {
		return 0, "", err
	}
	if s.checkContestStarting(contest) != nil {
		return 0, "", err
	}
	if s.checkContestFinished(contest) != nil {
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
	if s.checkTaskExist(task) != nil {
		return 0, "", err
	}
	if s.checkTaskNotStarted(teamTask) != nil {
		return 0, "", err
	}
	if s.checkTaskNextHintNumExist(teamTask, hintNum) != nil {
		return 0, "", err
	}
	// если номер подсказки для показа меньше чем текущий, то показываем последнюю показанную подсказку
	if s.isTaskHintAlreadyShown(teamTask, hintNum) {
		return teamTask.NextHintNum, getHintByNum(teamTask.NextHintNum-1, task), nil
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

	return teamTask.NextHintNum, getHintByNum(teamTask.NextHintNum, task), nil
}

func getHintByNum(num int, task taskEntity) string {
	if num <= 0 {
		return task.Hints[0]
	}
	if num >= len(task.Hints)-1 {
		return task.Hints[len(task.Hints)-1]
	}

	return task.Hints[num]
}
