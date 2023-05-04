package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
)

func (s PgRepo) GetContest(teamId int) (api.GetContestResponse, error) {
	contest, err := s.getContestEntity(teamId)
	if err != nil {
		return api.GetContestResponse{}, err
	}
	// @todo: pros проверить ошибки

	return api.GetContestResponse{
		ContestId: contest.Id,
		StartAt:   contest.StartAt,
		EndAt:     contest.EndAt,
	}, nil
}
