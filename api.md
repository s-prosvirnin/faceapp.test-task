# API

### Аутентификация

`POST /auth`  
`Response statuses: 200, 500`

- Экраны
    - Аутентификация
- Request
    
    ```jsx
    {
    	"login": "abc",
    	"pass": "abc"
    }
    ```
    
- Success response
    
    ```jsx
    {
    	"success": {
    		"auth_token": "abc",
    		"team_id": 123
    	}
    }
    ```
    
- Error response
    
    ```jsx
    {
    	"errors": {
    		"keys": [
    			"internal_error", // серверная ошибка
    			"login_pass_invalid", // неверный логин или пароль
    			"invalid_request" // поля в реквесте невалидны, см. invalid_request_fields
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["login", "pass"]
    	}
    }
    ```
    

### Данные по текущему турниру

`POST /team/contest`  
`Response statuses: 200, 500`

- Описание
    
    Если пользователь залогинился до начала турнира, нужно показывать какой-то ожидающий экран, чтобы список заданий все пользователи увидели в одно время. Ближе к концу турнира можно вызвать метод для синхронизации завершения турнира.
    
- Экраны
    - Аутентификация
        - если турнир не начался: можно показывать экран со временем отсчета до начала турнира
    - Список всех заданий команды
        - если турнир закончился
            - скрытие поля для ввода ответа по заданию
            - скрытие поля для получения подсказки
    - Отсчет до начала турнира (как вариант)
- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {
    	"team_id": 123
    }
    ```
    
- Success response
    
    ```jsx
    {
    	success: {
    		"contest_id": 123
    		"start_at": "1994-11-05T08:15:30Z"
    		"end_at": "1994-11-05T08:15:30Z"
    		"time_to_start_sec": 123
    		// статус турнира - турнир идет|скоро начнется|завершен
    		"status": "started|will_start_soon|completed"
    	}
    }
    ```
    
- Error Response
    
    ```jsx
    {
    	"errors": {
    		"keys": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"team_not_found" // команда не найдена
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
    

### Список всех заданий команды

`POST /team/contest/tasks`  
`Response statuses: 200, 500`

- Экраны
    - Авторизация
- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {
    	"team_id": 123
    }
    ```
    
- Success response
    
    ```jsx
    {
    	success: {
    		tasks: [
    			"id": 123,
    			"name": "abc",
    			"coords": {"lat": 123.456, "lon": 123.456},
    			"description": "abc",
    			// сдано|не начато|попытка сдачи провалена
    			"status": "passed|not_started|attempt_failed",
    			// отправленные ответы (этого поля нет в ТЗ, но, кажется, оно было бы полезным)
    			"answers": [{"answer": "abc", "is_passed": true}]
    			// подсказки
    			"hints": {
    				// открытые подсказки. номер подсказки - ключ массива
    				"opened": ["abc"],
    				// количество доступных подсказок
    				"total": 123,
    				// следующий номер подсказки. -1 - если все подсказки исчерпаны
    				// с этим номером можно идти в /team/contest/task/hint и получать следующую подсказку
    				"next_num": 123,
    			},
    		]
    	}
    }
    ```
    
- Error Response
    
    ```jsx
    {
    	"errors": {
    		"keys": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"team_not_found", // команда не найдена
    			"contest_not_found", // нет турниров
    			"contest_not_started" // турнир не начался
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
    

### Начать задание

`POST /team/contest/task/start`  
`Response statuses: 200, 500`

- Экраны
    - Список всех заданий команды
- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {
    	"team_id": 123,
    	"task_id": 123
    }
    ```
    
- Success response
    
    ```jsx
    {
    	success: true // дефолтное значение, всегда true
    }
    ```
    
- Error response
    
    ```jsx
    {
    	"errors": {
    		"keys": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"team_not_found", // команда не найдена
    			"contest_not_found", // нет турниров
    			"contest_not_started" // турнир не начался
    			"contest_finished" // турнир закончился
    			"task_not_found" // задание не найдено
    			"task_already_started" // задание уже начато
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
    

### Отправить ответ по заданию

`POST /team/contest/task/answer`  
`Response statuses: 200, 500`

- Экраны
    - Список всех заданий команды
- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {
    	"team_id": 123,
    	"task_id": 123,
    	"answer": "abc"
    }
    ```
    
- Success response
    
    ```jsx
    {
    	success: {
    		"answer_passed": true
    	}
    }
    ```
    
- Error response
    
    ```jsx
    {
    	"errors": {
    		"keys": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"team_not_found", // команда не найдена
    			"contest_not_found", // нет турниров
    			"contest_not_started" // турнир не начался
    			"contest_finished" // турнир закончился
    			"task_not_found" // задание не найдено
    			"answer_already_passed" // ответ по заданию уже принят
    			"answer_per_time_limit_exceeded" // превышено количество ответов за единицу времени
    			"answer_limit_exceeded" // превышено общее количество ответов
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
    

### Показать подсказку по заданию

`POST /team/contest/task/hint`  
`Response statuses: 200, 500`

- Экраны
    - Список всех заданий команды
- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {
    	"team_id": 123,
    	"task_id": 123,
    	"hint_num": 123 // номер подсказки, которую надо показать (0|1|2)
    }
    ```
    
- Success response
    
    ```jsx
    {
    	success: {
    		"hint": "abc",
    		// номер следующей доступной подсказки (0|1|2). -1 - если подсказки исчерпаны
    		"next_num": 123
    	}
    }
    ```
    
- Error response
    
    ```jsx
    {
    	"errors": {
    		"key": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"team_not_found", // команда не найдена
    			"contest_not_found", // нет турниров
    			"contest_not_started" // турнир не начался
    			"contest_finished" // турнир закончился
    			"task_not_found" // задание не найдено
    			"hint_num_not_exist" // нет подсказки с таким номером
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
    

### Результаты всех команд

`POST /contest/results`  
`Response statuses: 200, 500`

- Headers
    - `Authorization: Bearer auth_token`
- Request
    
    ```jsx
    {}
    ```
    
- Success response
    
    ```jsx
    {
    	success: {
    		"teams_results": [{
    			// порядковый номер команды в рейтинге
    			"team_rank": 123,
    			"team_name": "abc",
    			"tasks_results": [{
    					// остальную инфу по заданию можно взять из запроса /team/contest/tasks
    					"task_id": 123,
    					// сдано|не начато|попытка сдачи провалена
    					"status": "passed|not_started|attempt_failed",
    					"hints_opened_count": 123
    				}
    			],
    			// количество сданных заданий
    			"tasks_passed_count": 123,
    			// суммарное штрафное время, сек
    			"penalty_time_sec": 123
    		}]
    	}
    }
    ```
    
- Error response
    
    ```jsx
    {
    	"errors": {
    		"key": [
    			"internal_error", // серверная ошибка
    			"auth_token_invalid", // аутентификационный токен не валиден
    			"invalid_request", // поля в реквесте невалидны, см. invalid_request_fields
    			"contest_not_found", // нет турниров
    			"contest_not_started" // турнир не начался
    		],
    		// поля, не прошедшие валидацию
    		"invalid_request_fields": ["team_id"]
    	}
    }
    ```
