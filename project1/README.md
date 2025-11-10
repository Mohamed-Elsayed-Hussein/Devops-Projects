# Three-Tier App with Go, MySQL, and Nginx (HTTPS) using Docker Compose

This repository contains a Go backend, a MySQL database, and an Nginx reverse proxy configured for HTTPS, all orchestrated with Docker Compose.

Below is the detailed step-by-step workflow of what was done and how to run it.

---

## Step-by-Step Workflow

### 1. Prepare the Project Structure

Organized the project with `backend` and `database` folders:

```
project1/
├── backend
│   ├── Dockerfile
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── database
│   ├── env          # contains MYSQL_ROOT_PASSWORD
│   └── db-password  # password file for backend
├── nginx
│   ├── Dockerfile
│   ├── nginx.conf
│   └── generate-ssl.sh
└── compose.yaml
```

- The Go backend reads the database password from `/run/secrets/db-password`.

### 2. Set Up Docker Network (initial manual step)

Created a custom Docker network for backend & database communication:

```
docker network create   --driver bridge   --subnet 10.0.0.0/24   --gateway 10.0.0.1   go-backend
```

### 3. Run MySQL Container (initial manual step)

Started MySQL with the environment file and connected it to the network:
Notes: 
* the database should have example and the container name should be example this from source code .
```
docker run --rm  -d --network go-backend --network-alias db --name  example -p 3306:3306 --hostname db --env-file database/env mysql:8
```

### 4. Build and Run Go Backend Container (initial manual step)

- Built backend Docker image using multi-stage Dockerfile.
- Ran container, mounting the password file and connecting to `go-backend` network:

Notes: 
* the password shouldnot have any newline or spaces .

```
docker run  --rm  -p 8000:8000 -v ./backend/db-password:/run/secrets/db-password  --network go-backend --name backend-go backend-go:v1.0
```

- Backend automatically creates the database and `blog` table with initial data.

### 5. Add Nginx Reverse Proxy

- Created `nginx` folder with:
  - `Dockerfile` for building proxy image.
  - `nginx.conf` with upstream configuration to backend.
  - `generate-ssl.sh` to generate self-signed certificates.

- Configured proxy to:
  - Listen on HTTPS (443)
  - Forward requests to backend service (`backend:8000`)
  - Set required proxy headers.

### 6. Convert to Docker Compose

Created `compose.yaml` (docker-compose) with three services.

- Defined two networks:
  - `go-backend` for backend & database
  - `reverse-proxy` for connecting proxy externally

### 7. Build and Run All Services

Used Docker Compose with `--build` to ensure images rebuild if changed:

```
docker compose -f compose.yaml up -d --build
```

### 8. Verify Services

Check backend via Nginx HTTPS proxy:

```
curl -k https://localhost:8080/
# or 
curl http://localhost:8081
```

- Output: JSON list of blog posts.
- HTTP requests to HTTPS port return `400 Bad Request` (expected due to HTTPS).
- HTTP requests to HTTP port will return `301 Moved Permanently` .

### 9. Notes

- Backend automatically initializes the database on startup.
- Nginx handles HTTPS with self-signed certificates.
- Docker Compose orchestrates the three-tier application with networks and port mappings.

---

## Quick Start

1) Ensure Docker Engine and Docker Compose plugin are installed.
2) From the `project1` directory, run:

```
docker compose up -d --build
```

3) Test via proxy:

```
curl -k https://localhost:8080/
```

---

## Cleanup

```
docker compose down -v
```

This removes containers, networks, and volumes created by Compose.