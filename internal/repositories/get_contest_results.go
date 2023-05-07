package repositories

import (
	"database/sql"
	"sort"
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type teamsTasksEntity struct {
	TeamId               int            `db:"team_id"`
	TeamName             string         `db:"team_name"`
	TaskId               int            `db:"task_id"`
	TaskStatus           sql.NullString `db:"task_status"`
	TaskAnswersCreatedAt pq.StringArray `db:"task_answers_created_at"`
	TaskNextHitNum       sql.NullInt32  `db:"task_next_hit_num"`
}

const (
	// 15 минут штрафа за каждую взятую подсказку
	hintPenaltySec = 15 * 60
	// 30 минут штрафа за каждую неверную попытку сдачи
	incorrectAnswerPenaltySec = 30 * 60
)

func (r PgRepo) GetContestResults(teamId int) (api.GetContestResultsResponse, error) {
	contest, err := r.getContestEntity(teamId)
	if err != nil {
		return api.GetContestResultsResponse{}, err
	}
	if err = checkContestExist(contest); err != nil {
		return api.GetContestResultsResponse{}, err
	}
	if err = checkContestStarting(contest); err != nil {
		return api.GetContestResultsResponse{}, err
	}

	var teamsTasks []teamsTasksEntity
	query := `
		select
		    teamEntity.id as team_id
		    , teamEntity.name as team_name
		    , task.id as task_id
		    , team_task.status as task_status
		    , team_task.answers_created_at as task_answers_created_at
		    , team_task.next_hint_num as task_next_hit_num
		from contest_task
		inner join task on task.id = contest_task.task_id
		inner join contest_team on contest_team.contest_id = contest_task.contest_id
		inner join teamEntity on teamEntity.id = contest_team.team_id
		left join team_task on team_task.task_id = contest_task.task_id and team_task.team_id = teamEntity.id
		where contest_task.contest_id = $1
		order by teamEntity.id
	`
	err = r.db.Select(&teamsTasks, query, contest.Id)
	if err == sql.ErrNoRows || teamsTasks == nil {
		// сверху мы проверили актуальность турнира, такая ситуация является внутренней ошибкой
		return api.GetContestResultsResponse{}, utils.NewErrWithType(
			errors.New("contest results not found"),
			api.ErrorTypeInternal,
		)
	}
	if err != nil {
		return api.GetContestResultsResponse{}, wrapInternalError(err, "db.Select")
	}

	return makeContestResultsResponse(teamsTasks, contest)
}

func makeContestResultsResponse(teamsTasks []teamsTasksEntity, contest contestEntity) (
	api.GetContestResultsResponse,
	error,
) {
	teamResponsesById := make(map[int]api.TeamResultResponse)
	var teamResponses []api.TeamResultResponse
	for _, teamTask := range teamsTasks {
		teamResponse, isExist := teamResponsesById[teamTask.TeamId]
		if !isExist {
			teamResponse.TeamName = teamTask.TeamName
		}
		status := getTaskStatusForResponse(teamTask.TaskStatus.String)
		if status == api.TaskStatusPassed {
			teamResponse.TasksPassedCount++
			passedAnswerTime, err := pq.ParseTimestamp(
				time.UTC,
				teamTask.TaskAnswersCreatedAt[len(teamTask.TaskAnswersCreatedAt)-1],
			)
			if err != nil {
				return api.GetContestResultsResponse{}, wrapInternalError(err, "answer.ParseTimestamp")
			}
			// штрафное время: прошедшее с начала турнира до сдачи этого задания
			teamResponse.PenaltyTimeSec += int(passedAnswerTime.UTC().Sub(contest.StartAt.UTC()).Seconds())
			// штрафное время: плюс 15 минут штрафа за каждую взятую подсказку
			teamResponse.PenaltyTimeSec += int(teamTask.TaskNextHitNum.Int32) * hintPenaltySec
			// штрафное время: плюс 30 минут штрафа за каждую неверную попытку сдачи
			teamResponse.PenaltyTimeSec += (len(teamTask.TaskAnswersCreatedAt) - 1) * incorrectAnswerPenaltySec
		}
		teamResponse.TaskResults = append(
			teamResponse.TaskResults, api.TaskResultResponse{
				TaskId:           teamTask.TaskId,
				Status:           status,
				HintsOpenedCount: int(teamTask.TaskNextHitNum.Int32),
			},
		)
		teamResponsesById[teamTask.TeamId] = teamResponse
	}

	for _, teamResponse := range teamResponsesById {
		teamResponses = append(teamResponses, teamResponse)
	}

	sort.Slice(
		teamResponses,
		func(i, j int) bool {
			if teamResponses[i].TasksPassedCount == teamResponses[j].TasksPassedCount {
				return teamResponses[i].PenaltyTimeSec < teamResponses[j].PenaltyTimeSec
			}

			return teamResponses[i].TasksPassedCount > teamResponses[j].TasksPassedCount
		},
	)

	for i := range teamResponses {
		teamResponses[i].TeamRank = i + 1
	}

	return api.GetContestResultsResponse{TeamResults: teamResponses}, nil
}
