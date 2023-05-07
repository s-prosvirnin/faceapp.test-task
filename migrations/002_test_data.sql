-- +goose Up
/****** contest *****/
insert into contest (start_at, end_at, status)
values (NOW(), NOW() + INTERVAL '10 DAY', 'actual');
insert into contest (start_at, end_at, status)
values (NOW() + INTERVAL '1 MINUTE', NOW() + INTERVAL '10 MINUTE', 'actual');
insert into contest (start_at, end_at, status)
values (NOW() + INTERVAL '1 MINUTE', NOW() + INTERVAL '10 MINUTE', 'actual');
insert into contest (start_at, end_at, status)
values (NOW() + INTERVAL '1 MINUTE', NOW() + INTERVAL '10 MINUTE', 'passed');
;
/****** contest *****/

/****** contest_team *****/
insert into contest_team (contest_id, team_id)
values (2, 1);
insert into contest_team (contest_id, team_id)
values (1, 1);
insert into contest_team (contest_id, team_id)
values (1, 2);
insert into contest_team (contest_id, team_id)
values (1, 3);
insert into contest_team (contest_id, team_id)
values (1, 4);
;
/****** contest_team *****/

/****** team *****/
insert into team (id, name, login, password) values (1, 'team1', 'login1', 'pass1');
insert into team (id, name, login, password) values (2, 'team2', 'login2', 'pass2');
insert into team (id, name, login, password) values (3, 'team3', 'login3', 'pass3');
insert into team (id, name, login, password) values (4, 'team4', 'login4', 'pass4');
/****** team *****/

/****** task *****/
insert into task (name, coords_lat, coords_lon, description, answer, hints)
values ('task1', 1.23, 3.21, 'description1', 'answer1', '{"hint1.1", "hint1.2", "hint1.3"}');
insert into task (name, coords_lat, coords_lon, description, answer, hints)
values ('task2', 1.56, 4.41, 'description2', 'answer2', '{"hint2.1", "hint2.2", "hint2.3"}');
insert into task (name, coords_lat, coords_lon, description, answer, hints)
values ('task3', 1.6, 4.9, 'description3', 'answer3', '{"hint3.1", "hint3.2", "hint3.3"}');
insert into task (name, coords_lat, coords_lon, description, answer, hints)
values ('task4', 1.4, 2.4, 'description4', 'answer4', '{"hint4.1", "hint4.2", "hint4.3"}');
;
/****** task *****/

/****** contest_task *****/
insert into contest_task (contest_id, task_id) values (1, 1);
insert into contest_task (contest_id, task_id) values (1, 2);
insert into contest_task (contest_id, task_id) values (1, 3);
insert into contest_task (contest_id, task_id) values (1, 4);
;
/****** contest_task *****/