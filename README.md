# Team Betting

This is an application where you can register teams (i.e. countries for
Eurovision Song Contest) and connect multiple devices to create a real time
analysis of the competition.

An administrator can add a competition and the teams and create a lobby where
multiple devices may connect and in realtime bet and analyse the competition.

## Migrations

Assuming MySQL is running in docker as per `docker-compose.yaml`.

```sh
goose --dir migrations/ mysql "betting:betting@tcp(127.0.0.1:3306)/betting?parseTime=true&charset=utf8mb4&collation=utf8mb4_bin" up
```
