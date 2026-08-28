CREATE TABLE IF NOT EXISTS projects (
 id TEXT PRIMARY KEY, organization TEXT NOT NULL, name TEXT NOT NULL, cost_center TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'active',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS clusters (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, provider TEXT NOT NULL, status TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS resource_groups (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), name TEXT NOT NULL, cluster_ids TEXT[] NOT NULL DEFAULT '{}', node_type TEXT NOT NULL DEFAULT '',
 cpu_quota INTEGER NOT NULL DEFAULT 0, memory_quota_mb INTEGER NOT NULL DEFAULT 0, gpu_quota INTEGER NOT NULL DEFAULT 0, concurrency INTEGER NOT NULL DEFAULT 1, budget_tag TEXT NOT NULL DEFAULT '',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(project_id,name), UNIQUE(project_id,id)
);
CREATE TABLE IF NOT EXISTS queues (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), resource_group_id TEXT NOT NULL REFERENCES resource_groups(id), name TEXT NOT NULL, priority INTEGER NOT NULL DEFAULT 0, weight INTEGER NOT NULL DEFAULT 1, concurrency INTEGER NOT NULL DEFAULT 1, max_duration_seconds INTEGER NOT NULL DEFAULT 0, scheduling_policy TEXT NOT NULL DEFAULT 'fifo',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(resource_group_id,name), UNIQUE(project_id,id)
);
CREATE TABLE IF NOT EXISTS asset_versions (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), asset_type TEXT NOT NULL, name TEXT NOT NULL, version TEXT NOT NULL, digest TEXT NOT NULL, checksum TEXT NOT NULL, uri TEXT NOT NULL, permissions JSONB NOT NULL DEFAULT '{}'::jsonb, source TEXT NOT NULL DEFAULT '',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(project_id,name,version), UNIQUE(project_id,digest)
);
CREATE TABLE IF NOT EXISTS training_jobs (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), name TEXT NOT NULL, template_version TEXT NOT NULL DEFAULT '', image_digest TEXT NOT NULL, code_version TEXT NOT NULL, dataset_version TEXT NOT NULL, resource_request JSONB NOT NULL DEFAULT '{}'::jsonb, queue_id TEXT NOT NULL, state TEXT NOT NULL, output_uri TEXT NOT NULL DEFAULT '',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), CHECK (state IN ('draft','pending_validation','queued','allocating','running','succeeded','failed','cancelled','stopped','timeout')), UNIQUE(project_id,name)
);
ALTER TABLE training_jobs ADD CONSTRAINT fk_training_queue_project FOREIGN KEY (project_id, queue_id) REFERENCES queues(project_id, id);
CREATE TABLE IF NOT EXISTS training_job_assets (
 job_id TEXT NOT NULL REFERENCES training_jobs(id), asset_version_id TEXT NOT NULL REFERENCES asset_versions(id), PRIMARY KEY(job_id,asset_version_id)
);
CREATE TABLE IF NOT EXISTS experiment_runs (
 id TEXT PRIMARY KEY, training_job_id TEXT NOT NULL REFERENCES training_jobs(id), parameters JSONB NOT NULL DEFAULT '{}'::jsonb, metrics JSONB NOT NULL DEFAULT '{}'::jsonb, logs_uri TEXT NOT NULL DEFAULT '', artifacts_uri TEXT NOT NULL DEFAULT '', best BOOLEAN NOT NULL DEFAULT false, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(training_job_id)
);
CREATE TABLE IF NOT EXISTS model_versions (
 id TEXT PRIMARY KEY, experiment_run_id TEXT NOT NULL REFERENCES experiment_runs(id), version TEXT NOT NULL, format TEXT NOT NULL, signature TEXT NOT NULL DEFAULT '', metrics_summary JSONB NOT NULL DEFAULT '{}'::jsonb, approval_status TEXT NOT NULL DEFAULT 'pending', lifecycle TEXT NOT NULL DEFAULT 'active', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS online_services (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), model_version_id TEXT NOT NULL REFERENCES model_versions(id), name TEXT NOT NULL, image_digest TEXT NOT NULL, resource_spec JSONB NOT NULL DEFAULT '{}'::jsonb, endpoint TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'stopped', version_policy TEXT NOT NULL DEFAULT 'fixed'
);
CREATE TABLE IF NOT EXISTS audit_events (
 id TEXT PRIMARY KEY, actor TEXT NOT NULL, project_id TEXT REFERENCES projects(id), resource_type TEXT NOT NULL, resource_id TEXT NOT NULL, action TEXT NOT NULL, before JSONB, after JSONB, request_id TEXT NOT NULL DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_jobs_project_state ON training_jobs(project_id,state);
CREATE INDEX IF NOT EXISTS idx_jobs_project_name ON training_jobs(project_id,name);
CREATE INDEX IF NOT EXISTS idx_assets_project_name ON asset_versions(project_id,name);
CREATE INDEX IF NOT EXISTS idx_audit_project_time ON audit_events(project_id,created_at);
