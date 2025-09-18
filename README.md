# Web-Applications



Запуск мигррации

docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres

migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5436/backand_GO?sslmode=disable" up