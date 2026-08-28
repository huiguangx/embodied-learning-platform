package db

import (
	"os"
	"strings"
	"testing"
)

func TestSeedCounts(t *testing.T) {
	sql, err := os.ReadFile("seed.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(sql)
	checks := map[string]int{
		"EIP Demo Project": 1,
		"Cloud GPU Cluster": 1,
		"IDC GPU Cluster": 1,
		"Default Cloud": 1,
		"Default IDC": 1,
		"'default'": 1,
		"'priority'": 1,
		"'idc-default'": 1,
		"training-base": 1,
		"demo-images": 1,
		"demo-labels": 1,
		"resnet50-demo": 1,
		"yolov8-demo": 1,
		"pending-demo": 1,
		"resnet50-online": 1,
		"seed-001": 1,
	}
	for needle, want := range checks {
		if got := strings.Count(text, needle); got < want {
			t.Errorf("seed missing %q", needle)
		}
	}
}

func TestMigrationDefinesRequiredTablesAndStateConstraint(t *testing.T) {
	path := "migrations/001_init.sql"
	sql, err := os.ReadFile(path)
	if err != nil { t.Fatal(err) }
	text := string(sql)
	for _, name := range []string{"projects", "clusters", "resource_groups", "queues", "asset_versions", "training_jobs", "experiment_runs", "model_versions", "online_services", "audit_events"} {
		if !strings.Contains(text, "CREATE TABLE IF NOT EXISTS "+name) { t.Errorf("missing table %s", name) }
	}
	if !strings.Contains(text, "pending_validation") || !strings.Contains(text, "idx_jobs_project_state") { t.Fatal("missing job state/index") }
}
