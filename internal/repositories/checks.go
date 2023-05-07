package repositories

import (
	"time"

	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/lib/pq"
)

func checkContestStarting(contest contestEntity) error {
	if contest.StartAt.After(time.Now()) {
		return utils.NewErrWithType(api.ErrContestNotStarted, api.ErrorTypeDomain)
	}

	return nil
}

func checkContestFinished(contest contestEntity) error {
	if contest.EndAt.Before(time.Now()) {
		return utils.NewErrWithType(api.ErrContestFinished, api.ErrorTypeDomain)
	}

	return nil
}

func checkContestExist(contest contestEntity) error {
	if contest.Id <= 0 {
		return utils.NewErrWithType(api.ErrContestNotFound, api.ErrorTypeDomain)
	}

	return nil
}

func checkTaskExist(task taskEntity) error {
	if task.Id <= 0 {
		return utils.NewErrWithType(api.ErrTaskNotFound, api.ErrorTypeDomain)
	}

	return nil
}

func checkTaskAlreadyStarted(task teamTaskEntity) error {
	if task.TaskId > 0 {
		return utils.NewErrWithType(api.ErrTaskAlreadyStarted, api.ErrorTypeDomain)
	}

	return nil
}

func checkTaskNotStarted(task teamTaskEntity) error {
	if task.TaskId <= 0 {
		return utils.NewErrWithType(api.ErrTaskNotStarted, api.ErrorTypeDomain)
	}

	return nil
}

func checkTaskNextHintNumExist(task taskEntity, teamTask teamTaskEntity, nextHintNum int) error {
	if nextHintNum >= len(task.Hints) ||
		nextHintNum > teamTask.NextHintNum ||
		nextHintNum < 0 {
		return utils.NewErrWithType(api.ErrTaskHintNumNotExist, api.ErrorTypeDomain)
	}

	return nil
}

func isTaskHintAlreadyShown(task teamTaskEntity, nextHintNum int) bool {
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

func checkTaskPassed(task taskEntity, teamTask teamTaskEntity) error {
	for _, teamAnswer := range teamTask.Answers {
		if teamAnswer == task.Answer {
			return utils.NewErrWithType(api.ErrTaskAnswerAlreadyPassed, api.ErrorTypeDomain)
		}
	}

	return nil
}

func isNextHintNumLast(task taskEntity, nextNum int) bool {
	if nextNum >= len(task.Hints) {
		return true
	}

	return false
}

func checkAnswerPerTimeLimitExceed(teamTask teamTaskEntity) error {
	for i := len(teamTask.AnswersCreatedAt) - 1; i >= 0; i-- {
		lastAnswerTime, _ := pq.ParseTimestamp(
			time.UTC,
			teamTask.AnswersCreatedAt[i],
		)

		if time.Now().Sub(lastAnswerTime).Seconds() <= 60 &&
			len(teamTask.AnswersCreatedAt)-i > answerPerMinLimit {
			return utils.NewErrWithType(api.ErrTaskAnswerPerTimeLimitExceeded, api.ErrorTypeDomain)
		}
	}

	return nil
}

func checkAnswersLimitExceed(teamTask teamTaskEntity) error {
	if len(teamTask.AnswersCreatedAt) > answersTotalLimit {
		return utils.NewErrWithType(api.ErrTaskAnswersLimitExceeded, api.ErrorTypeDomain)
	}

	return nil
}
