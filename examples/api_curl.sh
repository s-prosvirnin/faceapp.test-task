# auth
curl -X POST --location "http://localhost:8085/auth" \
    -H "Content-Type: application/json" \
    -d "{
          \"login\": \"login1\",
          \"pass\": \"pass1\"
        }"

# get contest
curl -X POST --location "http://localhost:8085/team/contest" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1
        }"

# get tasks
curl -X POST --location "http://localhost:8085/team/contest/tasks" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1
        }"

# start task
curl -X POST --location "http://localhost:8085/team/contest/task/start" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1,
          \"task_id\": 1
        }"

# add task answer
curl -X POST --location "http://localhost:8085/team/contest/task/answer" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1,
          \"task_id\": 1,
          \"answer\": \"answer1\",
          \"answer_uuid\": \"1f2h87\"
        }"

# show task hint
curl -X POST --location "http://localhost:8085/team/contest/task/hint" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1,
          \"task_id\": 1,
          \"hint_num\": 0
        }"

# get contest results
curl -X POST --location "http://localhost:8085/contest/results" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 114812dd-515b-4be8-9f5f-2fb391e9411f" \
    -d "{
          \"team_id\": 1
        }"