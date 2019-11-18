## Requirements
* GoLang 1.13+
* PostgreSQL 11+(с расширением pg_trgm)
* Make
* Docker
* docker-compose

## Development
* Сборка происходит через Makefile
* Запуск окружения `make dev-docker-compose-up` и `make dev-docker-compose-down`
* Добавление/удаление зависимостей `make tidy`
* Скачивание зависимостей `make vendor`
* Генерация кода `make generate`
* Запуск тестов `make test`
* Компиляция проекта `make build`
По умолчание компилируется под linux. 
Если нужно под другую платформу, то нужно переопределить `GOOS` переменную окружения.
Для примера: `GOOS=darwin make build` 
После компиляции бинарный файл будет лежать в папке `bin`

# Структура проекта
* assets - папка для различных файлов миграций, дампов и т.д.
* cmd - папка для команд с кот. можно запустить бинарник
* configs - папка для конфигов
* container - helper для DI
* definition - общие реализации контейнеров логгера, бд и т.д.
* deployments - папка для docker-compose файлов
* internal - код проекта

## Примеры запуска
# Запуск http сервера 
Для запуска http сервера нужно указать ключ `-c` с полным путем до конфига. Сервер запустится на `9090` порту.
`./bin/theboatscom http-server -c /Users/igor.tumanov/go/src/gitlab.com/igor.tumanov1/theboatscom/configs/config.json`

# Запуск миграций
Для запуска миграций `-c` и `-m`, где `-c` - путь до конфига, `-m` путь до папки с миграциями 
`./bin/theboatscom migrate up -c /Users/igor.tumanov/go/src/gitlab.com/igor.tumanov1/theboatscom/configs/config.json -m /Users/igor.tumanov/go/src/gitlab.com/igor.tumanov1/theboatscom/asstets/migrations`

Откат последне миграции
`./bin/theboatscom migrate down -c /Users/igor.tumanov/go/src/gitlab.com/igor.tumanov1/theboatscom/configs/config.json -m /Users/igor.tumanov/go/src/gitlab.com/igor.tumanov1/theboatscom/asstets/migrations`

Дамп с тестовыми данными загрузится автоматически.
В задании говорится, что должна быть команда для загрузки данных. Это можно сделать через миграцию.
Соответственно, если нужно перезалить данные, то можно откатить последнюю миграцию.
Это временное решение, но ничего не мешает пренести эти операции в отдельную консолную команду.

#Примеры запросов к API
Curl запрос на поиск
`curl -iv -X GET "http://127.0.0.1:9090/api/search?query=test"`

Curl запрос на автодополнение
`curl -iv -X GET "http://127.0.0.1:9090/api/autocomplete?query=test"`
