#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
if command -v go >/dev/null 2>&1; then
  (cd "$ROOT/backend" && go test ./...)
else
  echo "go not found; skipped backend tests" >&2
fi
if [ -d "$ROOT/frontend/node_modules" ]; then
  (cd "$ROOT/frontend" && npm run build)
else
  echo "frontend/node_modules missing; run npm install before frontend build" >&2
fi
docker compose -f "$ROOT/docker-compose.yml" config >/dev/null
echo "MVP checks complete"
