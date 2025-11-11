# Database (MySQL 8)

This folder contains configuration for the MySQL database used by the Go backend.

## Files

- `env` — Environment file loaded by the MySQL container. Must include at least:

```
MYSQL_ROOT_PASSWORD=<your-strong-password>
```

- `db-password` — Password file mounted into the backend container at `/run/secrets/db-password`.

## Run (Direct Docker)

```
docker run -d --network go-backend --name example \
  --hostname db --env-file database/env mysql:8
```

## With Docker Compose

Handled automatically by `compose.yaml` in the project root. To start all services:

```
docker compose up -d --build
```

## Notes

- The backend connects to MySQL using hostname `db` on the `go-backend` network.
- Ensure `env` is present with a valid `MYSQL_ROOT_PASSWORD` before running.
- Do not commit real production passwords to version control. Use local `.gitignore` for sensitive files when needed.