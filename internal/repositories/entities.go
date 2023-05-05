package repositories

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const contestStatusActual = "actual"

type contestEntity struct {
	Id      int       `db:"id"`
	Status  string    `db:"status"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
}

type taskEntity struct {
	Id          int            `db:"id"`
	Name        string         `db:"name"`
	CoordsLat   float32        `db:"coords_lat"`
	CoordsLon   float32        `db:"coords_lon"`
	Description string         `db:"description"`
	Answer      string         `db:"answer"`
	Hints       pq.StringArray `db:"hints"`
}

type teamTaskEntity struct {
	TeamId           int            `db:"team_id"`
	TaskId           int            `db:"task_id"`
	Answers          pq.StringArray `db:"answers"`
	AnswersUuid      pq.StringArray `db:"answer_uuids"`
	AnswersCreatedAt pq.StringArray `db:"answers_created_at"`
	NextHintNum      int            `db:"next_hint_num"`
	Status           string         `db:"status"`
}

type team struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Login       string `db:"login"`
	Password    string `db:"password"`
	AccessToken string `db:"access_token"`
}

func (r PgRepo) getContestEntity(teamId int) (contestEntity, error) {
	var contest contestEntity
	query := `
		select c.*
		from contest c
		inner join contest_team ct on c.id = ct.contest_id
		where ct.team_id = $1 and c.status = $2
	`
	err := r.db.Get(&contest, query, teamId, contestStatusActual)
	if err == sql.ErrNoRows {
		return contestEntity{}, nil
	}
	if err != nil {
		return contest, wrapInternalError(err, "db.Get")
	}

	return contest, nil
}

func (r PgRepo) getTaskEntity(contestId int, taskId int) (taskEntity, error) {
	var task taskEntity
	query := `
		select
		    t.*
		from contest_task ct
		inner join task t on t.id = ct.task_id
		where ct.contest_id = $1 and  ct.task_id = $2
	`
	err := r.db.Get(&task, query, contestId, taskId)
	if err == sql.ErrNoRows {
		return taskEntity{}, nil
	}
	if err != nil {
		return taskEntity{}, wrapInternalError(err, "db.Get")
	}

	return task, nil
}

func (r PgRepo) getTeamTaskEntity(teamId int, taskId int) (teamTaskEntity, error) {
	var teamTask teamTaskEntity
	query := `
		select
		    tt.*
		from team_task tt
		where tt.team_id = $1 and tt.task_id = $2
	`
	err := r.db.Get(&teamTask, query, teamId, taskId)
	if err == sql.ErrNoRows {
		return teamTaskEntity{}, nil
	}
	if err != nil {
		return teamTaskEntity{}, wrapInternalError(err, "db.Get")
	}

	return teamTask, nil
}
