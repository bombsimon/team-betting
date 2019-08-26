# Team Betting

This is an application where you can register teams (i.e. countries for
Eurovision Song Contest) and connect multiple devices to create a real time
analysis of the competition.

An administrator can add a competition and the teams and create a lobby where
multiple devices may connect and in realtime bet and analyse the competition.

## Migrations

Assuming MySQL is running in docker as per `docker-compose.yaml`.

Install Goose `go get -u github.com/pressly/goose/cmd/goose` and run the
migration.

```sh
goose \
    --dir migrations/ \
    mysql \
    "betting:betting@/betting?parseTime=true&charset=utf8mb4&collation=utf8mb4_bin" \
    up
```

Or if you want to use Gorm to run migrations.

```sh
[ADD_DATA=1] [GET_DATA=1] go run cmd/gorm-migrate/main.go
```
