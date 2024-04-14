# avito_test
## Описание
Тестовый проект написан на Go.  
База для хранения основных данных Postgres.  
База данных для кеша Redis 

Сервис разворачивается через docker-compose up.  
Для корректной работы необходимо создать файл [.env]и прописать там настройки для подключения

    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=postgres


запуск миграции объектов базы
go run ./pkg/postgres_db/migration/main.go

Заполнение базы тестовыми данными
go run ./tests/test_data/main.go


запуск тестов
go test -v ./tests/ 

go run ./request_examples/main.go