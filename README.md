# rullafy-data-api

```bash
export GOPRIVATE=github.com/vanern/*
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