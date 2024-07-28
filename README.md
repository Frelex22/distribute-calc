## Структура проекта

awesomeProject3/
├── orchestrator/
│ ├── main.go
│ ├── handlers.go
│ ├── models.go
│ ├── Dockerfile
├── agent/
│ ├── main.go
│ ├── grpc_client.go
│ ├── Dockerfile
├── proto/
│ ├── task.proto
├── grpc-server/
│ ├── main.go
│ ├── Dockerfile
├── docker-compose.yml
├── README.md



## Установка и запуск

### Требования

- Docker
- Docker Compose

### Шаги для запуска

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/awesomeProject3.git
    cd awesomeProject3
    ```

2. Сгенерируйте gRPC код:
    ```sh
    protoc --go_out=plugins=grpc:. proto/task.proto
    ```

3. Запустите проект с помощью Docker Compose:
    ```sh
    docker-compose up --build
    ```

Теперь проект доступен на `http://localhost:8080` (для API оркестратора) и `http://localhost:50051` (для gRPC сервера).


Пример `curl` запроса:
```sh
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
    "login": "user1",
    "password": "password"
    
Авторизация пользователя
POST /api/v1/login
Content-Type: application/json

{
    "login": "user1",
    "password": "password"
}

Добавление арифметического выражения

POST /api/v1/calculate
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
    "expression": "2 + 2 * 2"
}

Получение выражения по его ID

GET /api/v1/expressions/{id}
Authorization: Bearer <JWT_TOKEN>