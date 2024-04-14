# avito_test
## Описание
Тестовый проект написан на Go.  
База для хранения основных данных Postgres.  
База данных для кеша Redis. 

Сервис разворачивается через docker-compose up.  
Для корректной работы необходимо создать файл [.env]и прописать там настройки для подключения

    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=postgres
    POSTGRES_HOST=postgres
    POSTGRES_HOST_LOCAL=localhost
    POSTGRES_PGDATA=/data/postgres
    REDIS_PASSWORD=sOmE_sEcUrE_pAsS
    REDIS_HOST=localhost
    REDIS_PORT=6379
    SECRET_KEY_TOKEN=SECRET_KEY_TOKEN

Запуск проекта:

    docker-compose up
Запуск миграции объектов базы:
```go
go run ./pkg/postgres_db/migration/main.go
```
Заполнение базы тестовыми данными:
```go
go run ./tests/test_data/main.go
```
Запуск теста:
```go
go test -v ./tests/ 
```
Запуск примеров запросов всех ручек:
```go
go run ./request_examples/main.go
```

TODO  
Не успел привести структуру кода в порядок.
Не успел покрыть тестами остальные ручки.
Начал писать функциональность для версионирования баннеров.  
Учел в структуре бд, но не успел дописать ручки.
Пришлось долго разбираться с аутентификацией, не сталкивался с ней.

