package repositories

func (r PgRepo) ShowTaskHint(teamId int, taskId int, hintNum int) (nextHintNum int, hint string, err error) {
	contest, err := r.getContestEntity(teamId)
	if err != nil {
		return 0, "", err
	}
	if err = checkContestExist(contest); err != nil {
		return 0, "", err
	}
	if err = checkContestStarting(contest); err != nil {
		return 0, "", err
	}
	if err = checkContestFinished(contest); err != nil {
		return 0, "", err
	}

	task, err := r.getTaskEntity(contest.Id, taskId)
	if err != nil {
		return 0, "", err
	}
	teamTask, err := r.getTeamTaskEntity(teamId, taskId)
	if err != nil {
		return 0, "", err
	}
	if err = checkTaskExist(task); err != nil {
		return 0, "", err
	}
	if err = checkTaskNotStarted(teamTask); err != nil {
		return 0, "", err
	}
	if err = checkTaskNextHintNumExist(task, teamTask, hintNum); err != nil {
		return 0, "", err
	}
	// если номер подсказки для показа меньше чем текущий, то показываем последнюю показанную подсказку
	if isTaskHintAlreadyShown(teamTask, hintNum) {
		return getHintNumForResponse(task, teamTask), getHint(task, teamTask.NextHintNum), nil
	}

	teamTask.NextHintNum++

	query := `
		update team_task
		set next_hint_num = :next_hint_num
		where team_id = :team_id and task_id = :task_id
	`
	res, err := r.db.NamedExec(query, teamTask)
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
