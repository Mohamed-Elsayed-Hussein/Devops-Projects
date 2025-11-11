# Backend (Go)

This service is a Go HTTP API that connects to MySQL. It automatically creates the database, initializes a `blog` table, and serves JSON responses.

## Key Details

- Secrets: The DB password is read from `/run/secrets/db-password` mounted at runtime.
- Port: The backend listens on port `8000` inside the container.
- DB Hostname: The MySQL hostname is `db` when using the provided Compose networks.
- Initialization: On startup, the service ensures the database and the `blog` table exist and seeds initial data.

## Build and Run (Direct Docker)

```
# Build image (example tag)
docker build -t backend-go:v1.0 .

# Run container on the custom network created earlier and mount the secret file
docker run -d --network go-backend \
  -v ./db-password:/run/secrets/db-password \
  --name backend-go backend-go:v1.0
```

## Build and Run (Compose)

From the `project1` directory:

```
docker compose up -d --build backend
```

## Endpoints

- `GET /` — returns a JSON array of blog posts.

Example test:

```
curl -k https://localhost:8080/
```

(When running via the Nginx proxy.)

## Configuration

- The application expects the following runtime dependencies:
  - MySQL server reachable at hostname `db` on the `go-backend` network.
  - Password file mounted at `/run/secrets/db-password`.

## Notes

- If you change code in `main.go`, re-run with `docker compose up -d --build` to rebuild the image.
- The backend is designed to create the database and table automatically on first run.