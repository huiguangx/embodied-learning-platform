# EIP Training Platform MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a runnable Next.js + Go Gin MVP that demonstrates the EIP training loop: dashboard, assets, resources, training jobs, experiments, model registration, and minimal online-service lifecycle.

**Architecture:** A single modular Gin API owns project context, RBAC, resources, assets, job state transitions, experiments, models, services, and audit events. A single Next.js app provides the console and calls only the Gin API. PostgreSQL is the source of truth; a local in-memory/mock scheduler adapter simulates execution without requiring Kubernetes.

**Tech Stack:** Next.js App Router, TypeScript, React, CSS modules or existing Next.js defaults, Go 1.22+, Gin, `database/sql`, PostgreSQL, Docker Compose, Go `testing`, and Playwright smoke tests.

**Spec:** `docs/eip-training-platform-prd-and-design.md`

## Global Constraints

- Do not connect to a real Kubernetes cluster, cloud API, registry, or object-storage control plane in the MVP.
- Keep external integrations behind interfaces with a deterministic local adapter.
- Every write API accepts `Idempotency-Key` and returns a `requestId`.
- Project authorization is enforced in the API; the browser must not contain privileged credentials.
- Training jobs reference immutable image digest, code version, and dataset version fields.
- Do not claim mock scheduler behavior is production scheduling.
- Use ASCII in source files unless the existing file already requires another encoding.
- Each non-trivial domain behavior gets a focused Go unit test; UI gets a smoke test for navigation and read-only lists.

## File Map

Create the following top-level structure:

```text
backend/
  cmd/server/main.go
  internal/config/config.go
  internal/db/{db.go,migrations/001_init.sql,seed.go}
  internal/auth/{middleware.go,roles.go}
  internal/domain/{models.go,errors.go}
  internal/repository/{projects.go,assets.go,resources.go,jobs.go,experiments.go,models.go,services.go,audit.go}
  internal/scheduler/{adapter.go,local.go}
  internal/http/{router.go,response.go,dashboard.go,assets.go,resources.go,jobs.go,experiments.go,models.go,services.go}
  internal/*_test.go
  go.mod
  Dockerfile
frontend/
  app/{layout.tsx,page.tsx,globals.css}
  app/(console)/{dashboard,development-machines,custom-training,experiments,model-repository,online-services,image-repository,object-storage,resource-groups}/page.tsx
  components/{shell,DataTable,StatusBadge,FilterBar}.tsx
  lib/{api,types,mock-data}.ts
  e2e/platform.spec.ts
  package.json
docker-compose.yml
README.md
```

### Task 1: Repository Bootstrap and Local Runtime

**Files:**
- Create: `backend/go.mod`, `backend/cmd/server/main.go`, `backend/internal/config/config.go`, `backend/Dockerfile`
- Create: `frontend/package.json`, `frontend/app/layout.tsx`, `frontend/app/globals.css`, `frontend/app/page.tsx`
- Create: `docker-compose.yml`, `README.md`

**Interfaces:**
- Produces `GET /healthz` returning `{ "status": "ok" }` from Gin.
- Produces a Next.js page with a stable shell and API base URL from `NEXT_PUBLIC_API_BASE_URL`.

- [ ] **Step 1: Write the failing backend health test**

```go
func TestHealthEndpoint(t *testing.T) {
    r := setupRouterForTest()
    req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
    rec := httptest.NewRecorder()
    r.ServeHTTP(rec, req)
    if rec.Code != http.StatusOK || rec.Body.String() != `{"status":"ok"}` {
        t.Fatalf("unexpected health response: %d %s", rec.Code, rec.Body.String())
    }
}
```

- [ ] **Step 2: Run `cd backend && go test ./...` and verify it fails because the router is absent.**
- [ ] **Step 3: Implement Gin bootstrap, config loading, `/healthz`, graceful shutdown, and a Next.js placeholder shell.**
- [ ] **Step 4: Run `cd backend && go test ./...`; run `cd frontend && npm run build`; verify both pass.**
- [ ] **Step 5: Commit with `git add backend frontend docker-compose.yml README.md && git commit -m "chore: bootstrap eip platform mvp"`.**

### Task 2: Database Schema, Seed Data, and Domain Types

**Files:**
- Create: `backend/internal/db/db.go`, `backend/internal/db/migrations/001_init.sql`, `backend/internal/db/seed.go`
- Create: `backend/internal/domain/models.go`, `backend/internal/domain/errors.go`
- Test: `backend/internal/db/db_test.go`, `backend/internal/domain/models_test.go`

**Interfaces:**
- Defines `Project`, `AssetVersion`, `ResourceGroup`, `Queue`, `TrainingJob`, `ExperimentRun`, `ModelVersion`, `OnlineService`, and `AuditEvent`.
- Defines job states `draft`, `pending_validation`, `queued`, `allocating`, `running`, `succeeded`, `failed`, `cancelled`, `stopped`, `timeout`.
- Seeds one project, two clusters, two resource groups, three queues, representative image/data assets, and sample jobs matching the screenshots.

- [ ] **Step 1: Write tests for valid state values and seed counts.**
- [ ] **Step 2: Run `cd backend && go test ./internal/domain ./internal/db`; verify failure before implementation.**
- [ ] **Step 3: Add PostgreSQL schema with UUID/string IDs, timestamps, project foreign keys, immutable asset references, indexes for project/status/name, and audit append-only table.**
- [ ] **Step 4: Implement seed insertion and domain validation helpers.**
- [ ] **Step 5: Run tests against a Docker PostgreSQL service and verify repeatable seeding.**
- [ ] **Step 6: Commit with `git add backend/internal/db backend/internal/domain && git commit -m "feat: add eip domain schema and seed data"`.**

### Task 3: Authentication Context, RBAC, and API Envelope

**Files:**
- Create: `backend/internal/auth/roles.go`, `backend/internal/auth/middleware.go`, `backend/internal/http/response.go`
- Test: `backend/internal/auth/middleware_test.go`, `backend/internal/http/response_test.go`

**Interfaces:**
- Reads development header `X-User-ID` and `X-Project-ID`; defaults to seeded demo user only when `APP_ENV=development`.
- Exposes `AuthContext{UserID, ProjectID, Role}` to handlers.
- JSON errors use `{ "code", "message", "requestId", "details" }`.

- [ ] **Step 1: Write tests for missing project context, allowed engineer reads, and denied non-admin writes.**
- [ ] **Step 2: Run `cd backend && go test ./internal/auth ./internal/http`; verify failure.**
- [ ] **Step 3: Implement middleware, roles (`engineer`, `data_engineer`, `operator`, `business_admin`), and request ID generation.**
- [ ] **Step 4: Run tests and verify unauthorized requests never reach repositories.**
- [ ] **Step 5: Commit with `git add backend/internal/auth backend/internal/http && git commit -m "feat: add project auth and api errors"`.**

### Task 4: Read APIs for Dashboard and Assets

**Files:**
- Create: `backend/internal/repository/projects.go`, `backend/internal/repository/assets.go`, `backend/internal/http/dashboard.go`, `backend/internal/http/assets.go`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/http/dashboard_test.go`, `backend/internal/http/assets_test.go`

**Interfaces:**
- `GET /api/v1/projects/:projectId/dashboard`
- `GET /api/v1/projects/:projectId/assets/images?namespace=&name=&page=`
- `GET /api/v1/projects/:projectId/assets/objects?bucket=&prefix=`
- All list endpoints return `{ items, nextCursor, total }` and enforce project scope.

- [ ] **Step 1: Write handler tests for dashboard metrics, image filters, object prefixes, and pagination.**
- [ ] **Step 2: Run focused tests and verify failure.**
- [ ] **Step 3: Implement repository queries and handlers using seeded data; return digest and checksum fields.**
- [ ] **Step 4: Run `cd backend && go test ./internal/http/...`; verify pass.**
- [ ] **Step 5: Commit with `git add backend/internal/repository backend/internal/http && git commit -m "feat: add dashboard and asset read apis"`.**

### Task 5: Resource Groups, Queues, and Reservations

**Files:**
- Create: `backend/internal/repository/resources.go`, `backend/internal/http/resources.go`
- Test: `backend/internal/repository/resources_test.go`, `backend/internal/http/resources_test.go`

**Interfaces:**
- `GET/POST /api/v1/projects/:projectId/resource-groups`
- `GET/POST /api/v1/projects/:projectId/queues`
- `GET/POST /api/v1/projects/:projectId/reservations`
- `POST /api/v1/projects/:projectId/resource-checks`
- Resource check returns `allowed`, `reasonCode`, `remaining`, and `alternatives`.

- [ ] **Step 1: Write tests for quota acceptance, quota rejection, reservation overlap, and role restrictions.**
- [ ] **Step 2: Run focused tests and verify failure.**
- [ ] **Step 3: Implement transactional quota checks and idempotent create handlers.**
- [ ] **Step 4: Run tests with concurrent quota requests and verify no over-allocation.**
- [ ] **Step 5: Commit with `git add backend/internal/repository/resources.go backend/internal/http/resources.go && git commit -m "feat: add resource governance apis"`.**

### Task 6: Training Job State Machine and Local Scheduler

**Files:**
- Create: `backend/internal/repository/jobs.go`, `backend/internal/scheduler/adapter.go`, `backend/internal/scheduler/local.go`, `backend/internal/http/jobs.go`
- Test: `backend/internal/domain/job_state_test.go`, `backend/internal/scheduler/local_test.go`, `backend/internal/http/jobs_test.go`

**Interfaces:**
- `POST /api/v1/projects/:projectId/training-jobs:validate`
- `GET/POST /api/v1/projects/:projectId/training-jobs`
- `POST /api/v1/projects/:projectId/training-jobs/:id:cancel`
- `type Scheduler interface { Submit(context.Context, TrainingJob) error; Cancel(context.Context, string) error }`
- `Transition(current, event) (JobState, error)` rejects illegal state regression.

- [ ] **Step 1: Write table-driven tests for every allowed transition and rejected regression.**
- [ ] **Step 2: Run `cd backend && go test ./internal/domain ./internal/scheduler ./internal/http`; verify failure.**
- [ ] **Step 3: Implement immutable asset-reference validation, quota precheck, idempotency-key storage, state transitions, and local adapter transitions from queued to running to succeeded.**
- [ ] **Step 4: Verify duplicate POST with the same key returns the original job and duplicate scheduler events do not duplicate transitions.**
- [ ] **Step 5: Commit with `git add backend/internal/domain backend/internal/scheduler backend/internal/repository/jobs.go backend/internal/http/jobs.go && git commit -m "feat: add training job lifecycle"`.**

### Task 7: Experiments, Models, Services, and Audit

**Files:**
- Create: `backend/internal/repository/experiments.go`, `backend/internal/repository/models.go`, `backend/internal/repository/services.go`, `backend/internal/repository/audit.go`
- Create: `backend/internal/http/experiments.go`, `backend/internal/http/models.go`, `backend/internal/http/services.go`
- Test: `backend/internal/http/experiments_test.go`, `backend/internal/http/models_test.go`, `backend/internal/http/services_test.go`

**Interfaces:**
- `GET /api/v1/projects/:projectId/experiment-runs`
- `POST /api/v1/projects/:projectId/model-versions`
- `GET/POST /api/v1/projects/:projectId/online-services`
- `POST /api/v1/projects/:projectId/online-services/:id:stop`
- Successful jobs automatically create one experiment run in the local adapter path.

- [ ] **Step 1: Write tests for automatic experiment creation, model provenance, publish validation, service stop, and audit records.**
- [ ] **Step 2: Run focused tests and verify failure.**
- [ ] **Step 3: Implement repositories and handlers; require a succeeded experiment for model registration and a published model for service creation.**
- [ ] **Step 4: Verify every write creates an audit event with actor, project, object, action, request ID, and before/after JSON.**
- [ ] **Step 5: Commit with `git add backend/internal/repository backend/internal/http && git commit -m "feat: add experiment model and service apis"`.**

### Task 8: Next.js Console Shell and Navigation

**Files:**
- Create: `frontend/components/shell/Sidebar.tsx`, `frontend/components/shell/Topbar.tsx`, `frontend/components/StatusBadge.tsx`, `frontend/lib/api.ts`, `frontend/lib/types.ts`
- Modify: `frontend/app/layout.tsx`, `frontend/app/page.tsx`, `frontend/app/globals.css`
- Test: `frontend/e2e/platform.spec.ts`

**Interfaces:**
- `api.get<T>(path, query)` calls `NEXT_PUBLIC_API_BASE_URL` and attaches demo project headers in development.
- Sidebar routes: `/`, `/development-machines`, `/custom-training`, `/experiments`, `/model-repository`, `/online-services`, `/image-repository`, `/object-storage`, `/resource-groups`.

- [ ] **Step 1: Write Playwright tests asserting shell navigation and active route labels.**
- [ ] **Step 2: Run `cd frontend && npm test` or the configured Playwright command and verify failure.**
- [ ] **Step 3: Implement screenshot-aligned shell: compact dark sidebar, top platform tabs, tabbed page area, tables, status badges, and responsive overflow.**
- [ ] **Step 4: Run the browser test and verify no horizontal overlap at 1280px and 390px widths.**
- [ ] **Step 5: Commit with `git add frontend && git commit -m "feat: add eip console shell"`.**

### Task 9: Dashboard, Development Machines, and Training Pages

**Files:**
- Create: `frontend/app/(console)/dashboard/page.tsx`, `frontend/app/(console)/development-machines/page.tsx`, `frontend/app/(console)/custom-training/page.tsx`
- Create: `frontend/components/DataTable.tsx`, `frontend/components/FilterBar.tsx`
- Modify: `frontend/lib/types.ts`
- Test: `frontend/e2e/training-pages.spec.ts`

**Interfaces:**
- Dashboard consumes `GET /dashboard` and renders Cloud/IDC tabs, cluster health, GPU/CPU/memory, task counts, and node details.
- Training page consumes `GET /training-jobs`, supports read-only filters, and shows copy/edit/delete controls disabled until mutation handlers are explicitly enabled.

- [ ] **Step 1: Write Playwright tests for dashboard metrics, empty development-machine state, training table columns, and pagination controls.**
- [ ] **Step 2: Run the tests and verify failure.**
- [ ] **Step 3: Implement pages using real API data with loading, empty, error, and stale-data states.**
- [ ] **Step 4: Verify page layout against the collected screenshots without triggering create/edit/delete actions.**
- [ ] **Step 5: Commit with `git add frontend/app frontend/components frontend/e2e && git commit -m "feat: add dashboard and training views"`.**

### Task 10: Asset, Resource, Experiment, Model, and Service Pages

**Files:**
- Create: `frontend/app/(console)/image-repository/page.tsx`, `frontend/app/(console)/object-storage/page.tsx`, `frontend/app/(console)/resource-groups/page.tsx`, `frontend/app/(console)/experiments/page.tsx`, `frontend/app/(console)/model-repository/page.tsx`, `frontend/app/(console)/online-services/page.tsx`
- Test: `frontend/e2e/asset-resource-pages.spec.ts`

**Interfaces:**
- Pages consume the read APIs from Tasks 4, 5, and 7.
- Experimental/model/service pages show the MVP workflow states and an explicit “development/mock adapter” label where no real external integration exists.

- [ ] **Step 1: Write Playwright tests for image filters, object-storage breadcrumb, resource quota cards, experiment rows, model provenance, and service status.**
- [ ] **Step 2: Run the tests and verify failure.**
- [ ] **Step 3: Implement tables, filters, breadcrumbs, resource cards, provenance links, and service status controls.**
- [ ] **Step 4: Verify keyboard focus, visible labels, empty/error states, and mobile overflow.**
- [ ] **Step 5: Commit with `git add frontend/app frontend/e2e && git commit -m "feat: add asset resource and delivery views"`.**

### Task 11: Integration, Docker Compose, and Verification

**Files:**
- Modify: `docker-compose.yml`, `README.md`
- Create: `scripts/check-mvp.sh`
- Test: `backend/internal/integration/api_test.go`, `frontend/e2e/smoke.spec.ts`

- [ ] **Step 1: Write an integration test that starts with seeded PostgreSQL, calls dashboard -> asset -> resource-check -> training validation, and asserts project isolation.**
- [ ] **Step 2: Run the integration test and verify failure before wiring services together.**
- [ ] **Step 3: Wire Compose services, migrations, seed command, Gin API, Next.js app, and health checks.**
- [ ] **Step 4: Implement `scripts/check-mvp.sh` to run Go tests, frontend lint/build, and Playwright smoke tests.**
- [ ] **Step 5: Run `docker compose up --build`, then `./scripts/check-mvp.sh`; verify all checks pass and document URLs, demo headers, and known mock limitations.**
- [ ] **Step 6: Commit with `git add . && git commit -m "test: verify eip platform mvp"`.**

## Definition of Done

- `docker compose up --build` starts PostgreSQL, Gin, and Next.js with seeded data.
- Dashboard, development machines, custom training, image repository, object storage, resource groups, experiments, model repository, and online services are navigable.
- Training validation rejects unauthorized assets and quota overages; duplicate idempotent requests do not duplicate jobs.
- A successful local training job creates an experiment run; a model can only be registered from a successful run; a service can only use a published model.
- Every write is project-scoped and audited.
- Go unit/integration tests, frontend build, and Playwright smoke tests pass.
- README clearly states that scheduler, registry, object storage, and serving are local adapters, not production integrations.

## Spec Coverage Check

- Current dashboard/development machine/training/image/object-storage pages: Tasks 4, 8, and 9.
- Resource groups, queues, reservations, quota governance: Task 5.
- Training state machine, immutable provenance, idempotency: Task 6.
- Experiments, model registry, online service lifecycle: Task 7 and 10.
- RBAC, audit, project isolation, security boundaries: Task 3 and 7.
- Observability and failure states: Tasks 4, 6, 9, and 11.
- External data/simulation integrations remain explicit future adapters per scope: Task 11 README limitations.
