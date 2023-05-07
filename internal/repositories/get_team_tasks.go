package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

func (r PgRepo) GetTeamTasks(teamId int) (api.GetTasksResponse, error) {
	contest, err := r.getContestEntity(teamId)
	if err != nil {
		return api.GetTasksResponse{}, err
	}
	if err = checkContestExist(contest); err != nil {
		return api.GetTasksResponse{}, err
	}
	if err = checkContestStarting(contest); err != nil {
		return api.GetTasksResponse{}, err
	}

	var tasks []taskEntity
	query := `
		select
		    t.*
		from contest_task ct
		inner join task t on t.id = ct.task_id
		where ct.contest_id = $1
	`
	err = r.db.Select(&tasks, query, contest.Id)
	if err == sql.ErrNoRows || tasks == nil {
		// сверху мы проверили актуальность турнира, такая ситуация является внутренней ошибкой
		return api.GetTasksResponse{}, utils.NewErrWithType(
			errors.New("team tasks not found"),
			api.ErrorTypeInternal,
		)
	}
	if err != nil {
		return api.GetTasksResponse{}, wrapInternalError(err, "tasks.db.Select")
	}

	var teamTasks []teamTaskEntity
	query = `
		select
		    tt.*
		from team_task tt
		where tt.team_id = $1
	`
	err = r.db.Select(&teamTasks, query, teamId)
	if err != nil {
		return api.GetTasksResponse{}, wrapInternalError(err, "teamTasks.db.Select")
	}

	return makeTasksResponse(tasks, teamTasks), nil
}

func makeTasksResponse(tasks []taskEntity, teamTasks []teamTaskEntity) api.GetTasksResponse {
	tasksResponse := make([]api.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		var teamTask teamTaskEntity
		for _, teamTaskTmp := range teamTasks {
			if teamTaskTmp.TaskId == task.Id {
				teamTask = teamTaskTmp
				break
			}
		}
		answersResponse := make([]api.AnswerResponse, 0, len(teamTask.Answers))
		for _, answer := range teamTask.Answers {
			answersResponse = append(
				answersResponse, api.AnswerResponse{
					Answer:   answer,
					IsPassed: answer == task.Answer,
				},
			)
		}
		nextHintNum := 0
		openedHints := []string{}
		status := getTaskStatusForResponse(teamTask.Status)
		if teamTask.TaskId > 0 {
			nextHintNum = getHintNumForResponse(task, teamTask)
			if isNextHintNumLast(task, teamTask.NextHintNum) {
				openedHints = task.Hints
			} else {
				openedHints = task.Hints[0:teamTask.NextHintNum]
			}
		}
		hintsResponse := api.HintsResponse{
			Total:   len(task.Hints),
			NextNum: nextHintNum,
			Opened:  openedHints,
		}
		tasksResponse = append(
			tasksResponse, api.TaskResponse{
				Id:   task.Id,
				Name: task.Name,
				Coords: api.CoordsResponse{
					Lat: task.CoordsLat,
					Lon: task.CoordsLon,
				},
				Description: task.Description,
				Status:      status,
				Answers:     answersResponse,
				Hints:       hintsResponse,
			},
		)
	}

	return api.GetTasksResponse{Tasks: tasksResponse}
}
