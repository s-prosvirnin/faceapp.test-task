package repositories

import (
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
)

func (r PgRepo) GetContest(teamId int) (api.GetContestResponse, error) {
	contest, err := r.getContestEntity(teamId)
	if err != nil {
		return api.GetContestResponse{}, err
	}
	if err = checkContestExist(contest); err != nil {
		return api.GetContestResponse{}, err
	}

	return makeGetContestResponse(contest), nil
}

func makeGetContestResponse(contest contestEntity) api.GetContestResponse {
	timeToStartSec := int(time.Now().Sub(contest.StartAt).Seconds())
	status := api.ContestStatusStarted
	if timeToStartSec < 0 {
		timeToStartSec = 0
	}

	if err := checkContestStarting(contest); err != nil {
		status = api.ContestStatusWillStartSoon
	}

	if err := checkContestFinished(contest); err != nil {
		status = api.ContestStatusCompleted
	}

	return api.GetContestResponse{
		ContestId:      contest.Id,
		StartAt:        contest.StartAt,
		EndAt:          contest.EndAt,
		TimeToStartSec: timeToStartSec,
		Status:         status,
	}
}
