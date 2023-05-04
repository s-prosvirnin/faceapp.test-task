/****** contest *****/
DROP TABLE IF EXISTS contest;
CREATE TABLE contest
(
    id BIGSERIAL PRIMARY KEY,
    status varchar,
    start_at timestamp,
    end_at timestamp
);
comment on table contest is 'Турниры';
comment on column contest.id is 'Идентификатор турнира';
comment on column contest.status is 'Текстовый статус актуальности турнира';
comment on column contest.start_at is 'Дата-время начала';
comment on column contest.end_at is 'Дата-время окончания';
/****** contest *****/

/****** contest_team *****/
DROP TABLE IF EXISTS contest_team;
CREATE TABLE contest_team
(
    contest_id BIGINT,
    team_id BIGINT,
    PRIMARY KEY(contest_id, team_id)
);
comment on table contest_team is 'Турниры команд';
comment on column contest_team.contest_id is 'Идентификатор турнира';
comment on column contest_team.team_id is 'Идентификатор команды';
/****** contest_team *****/

/****** team *****/
DROP TABLE IF EXISTS team;
CREATE TABLE team
(
    -- сделал упрощение, что токен один на команду, по-правильному, нужно привязывать к пользователю (возможно, к устройству)
    id BIGINT PRIMARY KEY,
    name VARCHAR,
    login VARCHAR,
    password VARCHAR,
    access_token VARCHAR
);
comment on table team is 'Команды';
comment on column team.id is 'Идентификатор команды';
comment on column team.name is 'Название';
comment on column team.login is 'Логин';
comment on column team.password is 'Пароль';
comment on column team.access_token is 'Аутентификационный токен';
/****** team *****/

/****** task *****/
DROP TABLE IF EXISTS task;
CREATE TABLE task
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR,
    coords_lat double precision,
    coords_lon double precision,
    description varchar,
    answer varchar,
    hints varchar[]
);
comment on table task is 'Задания';
comment on column task.id is 'Идентификатор';
comment on column task.name is 'Название';
comment on column task.coords_lat is 'Координата задания: широта';
comment on column task.coords_lon is 'Координата задания: долгота';
comment on column task.description is 'Описание';
comment on column task.answer is 'Правильный ответ';
comment on column task.hints is 'Массив подсказок';
/****** task *****/

/****** contest_task *****/
DROP TABLE IF EXISTS contest_task;
CREATE TABLE contest_task
(
    contest_id BIGINT,
    task_id BIGINT,
    PRIMARY KEY(contest_id, task_id)
);
comment on table contest_task is 'Связь турниров и заданий';
comment on column contest_task.contest_id is 'Идентификатор турнира';
comment on column contest_task.task_id is 'Идентификатор задания';
/****** contest_task *****/

/****** team_task *****/
DROP TABLE IF EXISTS team_task;
CREATE TABLE team_task
(
    team_id BIGINT,
    task_id BIGINT,
    PRIMARY KEY(team_id, task_id),
    answers varchar[],
    answer_uuids varchar[],
    answers_created_at timestamp[],
    next_hint_num int,
    status varchar
);
comment on table team_task is 'Задания команды';
comment on column team_task.team_id is 'Идентификатор команды';
comment on column team_task.task_id is 'Идентификатор задания';
comment on column team_task.answers is 'Массив отправленных ответов';
comment on column team_task.answer_uuids is 'Массив uuid отправленных ответов';
comment on column team_task.answers_created_at is 'Массив дата-времени принятия отправленных ответов';
comment on column team_task.next_hint_num is 'Номер следующей неоткрытой подсказки. -1, если больше нет подсказок';
comment on column team_task.status is 'Текстовый статус';
/****** team_task *****/