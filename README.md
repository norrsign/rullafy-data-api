# rullafy-data-api





- pgx https://github.com/code-heim/go_74_pgx

```bash
export GOPRIVATE=github.com/vanern/*
```

To install PGGEN 
```bash
curl -L https://github.com/jschaf/pggen/releases/latest/download/pggen-linux-amd64.tar.xz \
  | tar -xJf - -C ~/bin

```

Run the Docker image 

```
# Build
docker build -t rullafy-data-api:latest .

# Basic run (public endpoints only)
docker run -p 8080:8080 goapi:latest

# With Keycloak integration (must supply the realm URL)
docker run -p 8080:8080 \
  -e GOAPI_REALM_URL=https://auth.dev.rullafy.techdevenv.com/realms/rullafy-dev \
  rullafy-data-api:latest

# Change port & realm in one go
docker run -p 9090:9090 \
  -e GOAPI_PORT=9090 \
  -e GOAPI_REALM_URL=https://auth.example.com/realms/prod \
  rullafy-data-api:latest

```
## Migration
- https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

```bash
curl -L  https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz  | tar xvz
 mv migrate  ~/go/bin/
migrate create -ext sql -dir db/migrations -seq create_authors_table

export DATABASE_URL="postgresql://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"
migrate -path db/migrations -database "$DATABASE_URL" up
migrate -path db/migrations -database "$DATABASE_URL" force 2

```