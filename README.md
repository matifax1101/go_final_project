# Планировщик задач

Веб-сервер на Go, реализующий TODO-планировщик с хранением задач в SQLite.

## Реализовано

- раздача фронтенда из `web/`;
- SQLite-база `scheduler.db`;
- добавление задач;
- получение списка ближайших задач;
- получение задачи по ID;
- редактирование задачи;
- удаление задачи;
- отметка задачи выполненной;
- повторение по правилам `d N` и `y`;
- обработчик `/api/nextdate`;
- переменные окружения `TODO_PORT` и `TODO_DBFILE`;
- аутентификация через `TODO_PASSWORD` и JWT-подобный подписанный токен с временем жизни 8 часов;
- Dockerfile для сборки и запуска приложения.

Не реализованы две дополнительные задачи со звёздочкой: правила повторения `w ...` и `m ...`, а также поиск задач через параметр `search`.

## Запуск локально

Требуется Go 1.24+.

```bash
go mod tidy
go run .
```

После запуска открыть:

`http://localhost:7540/`

Порт можно изменить:

```bash
TODO_PORT=8080 go run .
```

Путь к БД:

```bash
TODO_DBFILE=./data/scheduler.db go run .
```

Для включения аутентификации:

```bash
TODO_PASSWORD=12345 go run .
```

## Тестирование

Сначала запустить сервер в отдельном терминале:

```bash
go run .
```

Затем:

```bash
go test -run ^TestApp$ ./tests
go test -run ^TestDB$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
go test ./tests
```

В `tests/settings.go` базовые тесты используют:

```go
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = ``
```

`FullNextDate` и `Search` оставляются `false`, поскольку соответствующие дополнительные задания не реализуются.

## Docker

Сборка:

```bash
docker build -t go-final-project .
```

Для хранения SQLite-файла на хосте создайте каталог `data` и запустите:

```bash
docker run --rm -p 7540:7540 -v "${PWD}/data:/data" go-final-project
```

Windows PowerShell:

```powershell
docker run --rm -p 7540:7540 -v "${PWD}/data:/data" go-final-project
```

Планировщик будет доступен по адресу `http://localhost:7540/`, а база будет находиться на хосте в `data/scheduler.db`.

Для Docker с паролем:

```bash
docker run --rm -p 7540:7540 -e TODO_PASSWORD=12345 -v "${PWD}/data:/data" go-final-project
```
