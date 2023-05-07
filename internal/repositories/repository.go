package repositories

import (
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/api"
	"github.com/fa-rda/high-tech-cross.sergei-prosvirin/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PgRepo struct {
	db *sqlx.DB
}

func NewPgRepo(db *sqlx.DB) PgRepo {
	return PgRepo{db: db}
}

func wrapInternalError(err error, wrapMessage string) error {
	return utils.NewErrWithType(errors.Wrap(err, wrapMessage), api.ErrorTypeInternal)
}

func getHintNumForResponse(task taskEntity, teamTask teamTaskEntity) int {
	if isNextHintNumLast(task, teamTask.NextHintNum) {
		return -1
	}

	return teamTask.NextHintNum
}

func getTaskStatusForResponse(status string) string {
	// если задание не начато, то нет записи в teamTask - проставить дефолтный статус
	if status == "" {
		return api.TaskStatusNotStarted
	}

	return status
}
