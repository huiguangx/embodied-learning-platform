# EIP Training Platform MVP

Local bootstrap for the EIP training platform MVP.

## Run locally

```sh
docker compose up --build
```

- API: <http://localhost:8080>
- Health: <http://localhost:8080/healthz>
- Console: <http://localhost:3000>

The console reads `NEXT_PUBLIC_API_BASE_URL` and shows it in the shell.

## Development checks

```sh
cd backend && go test ./...
cd frontend && npm run build
```

Task 1 intentionally includes only local runtime bootstrap: Gin health,
configuration, graceful shutdown, a placeholder Next.js shell, and Compose.
