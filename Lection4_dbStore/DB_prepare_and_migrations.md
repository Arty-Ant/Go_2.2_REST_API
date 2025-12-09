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


# 2. Start

Создать директорию и файл  

- directory: `db/migrations`
- file : `docker-compose.yml`

Команты:
go mod init Bankstore
cd ..
go work use Bankstore/
cd Bankstore
mkdir -p db/migrations
touch docker-compose.yaml
touch justfile //команды для Just
touch .env //конфигурация
touch .env.example для отправки в репозиторий чтоб не попал по git ignore

### Folder and file structures:
```
bankstore/

❯ tree   
.
├── justfile
├── .env
├── db
│   └── migrations/
└── docker-compose.yml
```



## Запуск Postgres DB используя `docker compose`

`.env` файл: 

```
POSTGRES_USER=app_user
POSTGRES_PASSWORD=pswd
POSTGRES_DB=bankdb
```


`docker-compose.yml` файл: 

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
      test: ["CMD-SHELL", "pg_isready -U app_user -d ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

volumes:
  bankdb-data:
```

## Запуск `docker compose`
```
docker compose up -d
```


### Check container

```
docker ps
```

*Если возникнет ошибка доступа к сокету:*
- [либо добавить текущего пользователя в группу docker](https://www.digitalocean.com/community/questions/how-to-fix-docker-got-permission-denied-while-trying-to-connect-to-the-docker-daemon-socket)
- либо изменить права доступа `sudo chmod 666 /var/run/docker.sock`
- либо добавить `sudo` вначале каждой команды


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

[Add SQLC library](https://github.com/sqlc-dev/sqlc/blob/main/docs/overview/install.md)