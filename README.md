# Golang-OpenApi

## Requirements <a name="Requirements"></a>
### Software: <a name="Software"></a>
- Go 1.24.1 -> https://go.dev/dl/
- Docker 4.28.0 -> https://docs.docker.com/desktop/install/ubuntu/


### Environment variables: <a name="EnvironmentVariables"></a>
- DB_URL=localhost:5432/postgres
- DB_USER=postgres
- DB_PASS=mysecretpassword
- DB_ADAPTER=postgresql
- PORT=8080
- OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
- OTEL_SERVICE_NAME=mi-server-http

### Start PostgreSQL <a name="StartPostgreSQL"></a>
```bash
docker run --name golang-openapi-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:16.3-alpine3.18
```


### Start application <a name="StartApplication"></a>
```bash
DATABASE=postgres ./launch.sh
```

## Development <a name="development"></a>
### Visual Studio Code Extensions: <a name="vscode-extensions"></a>
#### Go

Install "Go" from go.dev: https://marketplace.visualstudio.com/items?itemName=golang.go


## Endpoints

| Name                 | Endpoint                                                             |
| -------------------- | -------------------------------------------------------------------- |
| Service Info         | http://localhost:8080/                                               |


## Observability <a name="observability"></a>

### Start Docker Grafana OTEL <a name="start-docker-grafana-otel"></a>

```bash
docker run --restart unless-stopped --detach --publish 3000:3000 --publish 4317:4317 --publish 4318:4318 --name grafana_otel grafana/otel-lgtm:0.11.0
```

Ref: https://hub.docker.com/r/grafana/otel-lgtm

### View Grafana <a name="view-grafana"></a>

Log in to http://localhost:3000 with user admin and password admin.