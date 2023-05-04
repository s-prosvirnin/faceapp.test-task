package repositories

import (
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
)

// @todo: pros убрать ресивер

func (s PgRepo) checkContestStarting(contest contestEntity) error {
	if contest.StartAt.After(time.Now()) {
		return utils.NewErrWithType(api.ErrContestNotStarted, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkContestFinished(contest contestEntity) error {
	if contest.EndAt.Before(time.Now()) {
		return utils.NewErrWithType(api.ErrContestFinished, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkContestExist(contest contestEntity) error {
	if contest.Id <= 0 {
		return utils.NewErrWithType(api.ErrContestNotFound, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkTaskExist(task taskEntity) error {
	if task.Id <= 0 {
		return utils.NewErrWithType(api.ErrTaskNotFound, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkTaskAlreadyStarted(task teamTaskEntity) error {
	if task.TaskId >= 0 {
		return utils.NewErrWithType(api.ErrTaskAlreadyStarted, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkTaskNotStarted(task teamTaskEntity) error {
	if task.TaskId <= 0 {
		return utils.NewErrWithType(api.ErrTaskNotStarted, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) checkTaskNextHintNumExist(task teamTaskEntity, nextHintNum int) error {
	if nextHintNum > task.NextHintNum || nextHintNum < 0 {
		return utils.NewErrWithType(api.ErrTaskHintNumNotExist, api.ErrorDomainType)
	}

	return nil
}

func (s PgRepo) isTaskHintAlreadyShown(task teamTaskEntity, nextHintNum int) bool {
	if nextHintNum < task.NextHintNum {
		return true
	}

	return false
}

func isAnswerPassed(task taskEntity, answer string) bool {
	if task.Answer == answer {
		return true
	}

	return false
}

func isTaskPassed(task taskEntity, teamTask teamTaskEntity) bool {
	for _, teamAnswer := range teamTask.Answers {
		if teamAnswer == task.Answer {
			return true
		}
	}

	return false
}
