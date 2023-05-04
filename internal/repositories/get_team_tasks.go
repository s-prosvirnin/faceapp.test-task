package repositories

import (
	"database/sql"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/pkg/errors"
)

func (s PgRepo) GetTeamTasks(teamId int) (api.GetTasksResponse, error) {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return api.GetTasksResponse{}, err
	}
	if s.checkContestExist(contest) != nil {
		return api.GetTasksResponse{}, err
	}
	if s.checkContestStarting(contest) != nil {
		return api.GetTasksResponse{}, err
	}
	if s.checkContestFinished(contest) != nil {
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
	err = s.db.Select(&tasks, query, contest.Id)
	if err == sql.ErrNoRows {
		// сверху мы проверили актуальность турнира, такая ситуация является внутренней ошибкой
		return api.GetTasksResponse{}, utils.NewErrWithType(
			errors.New("contest results not found"),
			api.ErrorInternalType,
		)
	}
	if err != nil {
		return api.GetTasksResponse{}, wrapInternalError(err, "db.Select")
	}

	var teamTasks []teamTaskEntity
	query = `
		select
		    tt.*
		from team_task tt
		where tt.team_id = $1 and tt.task_id = $2
	`
	err = s.db.Select(&teamTasks, query, teamId)
	if err != nil {
		return api.GetTasksResponse{}, wrapInternalError(err, "db.Select")
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
		status := api.TaskStatusNotStarted
		if teamTask.TaskId > 0 {
			nextHintNum = teamTask.NextHintNum
			openedHints = task.Hints
			if nextHintNum >= 0 {
				openedHints = task.Hints[0:teamTask.NextHintNum]
			}
			if teamTask.Status != "" {
				status = teamTask.Status
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
