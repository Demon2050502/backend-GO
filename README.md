# Web-Applications

Запуск мигррации

docker rm -f todo-db  
docker run --name=todo-db \
 -e POSTGRES_PASSWORD='1' \
 -e POSTGRES_DB='backand_GO' \
 -p 5436:5432 -d --rm postgres

migrate -path ./migrations -database "postgres://postgres:1@localhost:5436/backand_GO?sslmode=disable" up

go run cmd/main.go

http://localhost:8000/auth/sign-up
