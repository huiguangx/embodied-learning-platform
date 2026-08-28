# EIP Training Platform MVP

Local bootstrap for the EIP training platform MVP.

## Run locally

```sh
docker compose up --build
```

- API: <http://localhost:8080>
- Health: <http://localhost:8080/healthz>
- Console: <http://localhost:3000>

The console reads `NEXT_PUBLIC_API_BASE_URL`. The current scheduler, asset
stores, and serving controller are local MVP adapters; no cloud, Kubernetes,
registry, or object-storage control plane is contacted.

## Development checks

```sh
cd backend && go test ./...
cd frontend && npm run build
```

## Verification

```sh
./scripts/check-mvp.sh
```

The backend currently starts without a database connection; PostgreSQL
migrations and seed data are available for the next persistence wiring step.
