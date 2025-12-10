# 1. Before

**How To Install and Use Docker Compose on Ubuntu 24.04**

[Docker install manual](https://docs.docker.com/engine/install/ubuntu/)  

**Docker install через bash script**
```
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

```

**Add Just**  
[Just a command runner](https://github.com/casey/just)

```
sudo apt install just
```
Добавить расширение в VSCode:
`nefrob.vscode-just-syntax`


# 2. Start

Создать необходимые директории и файлы    

- directory: `db/migrations`
- file : `docker-compose.yml`
- file : `.env`
- file : `justfile`

**Команды:**
```bash
go mod init Bankstore
cd ..
go work use Bankstore/
cd Bankstore
mkdir -p db/migrations
touch docker-compose.yaml
touch justfile # команды для Just
touch .env # конфигурация для БД, этот файл не попадает в репозиторий
touch .env.example # Необязательно. Копия файла .env, для отправки в репозиторий
```

### Folder and file structures:
```
Bankstore/

❯ tree   
.
├── justfile
├── .env
├── db
│   └── migrations/
└── docker-compose.yaml
```



## Запуск Postgres DB используя `docker compose`

`.env` файл: 

```
POSTGRES_USER=app_user
POSTGRES_PASSWORD=pswd
POSTGRES_DB=bankdb
```


`docker-compose.yaml` файл: 

```
services:
  db:
    container_name: postgresdb
    image: postgres:17.2
    environment:
      POSTGRES_DB: "${POSTGRES_DB}"
      POSTGRES_USER: "${POSTGRES_USER}"
      POSTGRES_PASSWORD: "${POSTGRES_PASSWORD}"
    volumes:
      - bankdb-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app_user -d ${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

volumes:
  bankdb-data:
```

## Создание и Запуск `docker compose`
```
docker compose up -d
```
Для остановки использовать команду docker stop
А для запуска docker start

### Check container

```
docker ps
```

*Если возникнет ошибка доступа к сокету:*
- [либо добавить текущего пользователя в группу docker](https://www.digitalocean.com/community/questions/how-to-fix-docker-got-permission-denied-while-trying-to-connect-to-the-docker-daemon-socket)
- либо изменить права доступа `sudo chmod 666 /var/run/docker.sock`
- либо добавить `sudo` вначале каждой команды

### Настройка justfile для быстрых команд

```
set dotenv-load

# Start postgresql service
pg_start:
    docker compose start

# Stop postgresql service
pg_stop:
    docker compose stop
```

### Check database 
```
# Docker compose command
docker compose exec -it db psql -U app_user -d bankdb

# Check existing database
postgres=# \l
```

# 3. Golang-migrate CLI

## Install golang-migrate cli

```
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar -xvz -C /home/${USER}/go/bin
```
*директория ~/go/bin/ должна быть добавлена в `PATH`*

## Create migration files

*Шаблон:*  
`migrate create -ext sql -dir YOUR/DIR -seq MIGRATION_NAME`

Создадим 2 миграции:
```
# Command
migrate create -ext sql -dir db/migrations -seq create_table

# Output
/db/migrations/000001_create_table.up.sql
/db/migrations/000001_create_table.down.sql
```
### Files format
```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

### Create SQL

Создаем таблицы для БД


### Run migration

Шаблон:  
`migrate -database YOUR_DATABASE_URL -path PATH_TO_YOUR_MIGRATIONS up
`  

Накатить миграцию:
```
migrate -path db/migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/bankdb?sslmode=disable" -verbose up
```
### Check new table in DB
```
docker compose -it exec db psql -U app_user
```
```
postgres=# \c bankdb
```

Откатить миграцию:
```
migrate -path db/migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/bankdb?sslmode=disable" -verbose down
```

# 4. CRUD

Генератор кода для перевода кода с Go на SQL и наоборот
[Add SQLC library](https://github.com/sqlc-dev/sqlc/blob/main/docs/overview/install.md)
После установки надо создать в Bankstore файл sqlc.yaml скопривать из документации код и поправить его под свои настройки
Пример
```
version: "2"
cloud:
  # Replace <PROJECT_ID> with your project ID from the sqlc Cloud dashboard
  project: "<PROJECT_ID>"
sql:
  - engine: "postgresql"
    queries: "./db/queries"
    schema: "./db/migrations"
    gen:
      go:
        package: "db"
        out: "./db/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_exact_table_names: false
```
Подробнее в разделе конфигурации/go в документации на сайте SQLC

Создать папки в Bankstore/db - sqlc и queries.
В папке queries создать файл account.sql и прописать в нём код на языке SQL с флагами *-- name: CreateAccount :one* из документации SQLC

Пример
```
-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * 
FROM accounts
WHERE id=$1
LIMIT 1;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2; --смещение

-- name: UpdateAccount :one
UPDATE accounts
SET balance=$1
WHERE id=$2
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
```
Потом запустить команду 
sqlc generate

Её можно прописать в just
# Generate sqlc code
sqlc:
    sqlc generate

Потом в папке Bankstore запустить программу go mod tidy