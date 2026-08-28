package db

import (
	"database/sql"
	"fmt"
)

const (
	DemoProjectID = "00000000-0000-0000-0000-000000000001"
)

// Seed is repeatable: all records use stable IDs and upserts.
func Seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil { return err }
	defer tx.Rollback()
	statements := []string{
		`INSERT INTO projects (id,organization,name,cost_center,status) VALUES ('00000000-0000-0000-0000-000000000001','EIP','EIP Demo Project','CC-ML-001','active') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO clusters (id,name,provider,status) VALUES ('00000000-0000-0000-0000-000000000101','Cloud GPU Cluster','cloud','healthy'),('00000000-0000-0000-0000-000000000102','IDC GPU Cluster','idc','healthy') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO resource_groups (id,project_id,name,cluster_ids,node_type,cpu_quota,memory_quota_mb,gpu_quota,concurrency,budget_tag) VALUES ('00000000-0000-0000-0000-000000000201','00000000-0000-0000-0000-000000000001','Default Cloud','{00000000-0000-0000-0000-000000000101}','GPU','256','1048576','16','8','cloud'),('00000000-0000-0000-0000-000000000202','00000000-0000-0000-0000-000000000001','Default IDC','{00000000-0000-0000-0000-000000000102}','GPU','128','524288','8','4','idc') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO queues (id,resource_group_id,name,priority,weight,concurrency,max_duration_seconds,scheduling_policy) VALUES ('00000000-0000-0000-0000-000000000301','00000000-0000-0000-0000-000000000201','default',50,1,8,86400,'fifo'),('00000000-0000-0000-0000-000000000302','00000000-0000-0000-0000-000000000201','priority',100,2,4,172800,'priority'),('00000000-0000-0000-0000-000000000303','00000000-0000-0000-0000-000000000202','idc-default',50,1,4,86400,'fifo') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO asset_versions (id,project_id,asset_type,name,version,digest,checksum,uri,source) VALUES ('00000000-0000-0000-0000-000000000401','00000000-0000-0000-0000-000000000001','image','training-base','v1','sha256:demo-image-v1','sha256:demo-image-v1','registry://eip/training-base@sha256:demo-image-v1','seed'),('00000000-0000-0000-0000-000000000402','00000000-0000-0000-0000-000000000001','dataset','demo-images','v1','sha256:demo-dataset-v1','sha256:demo-dataset-v1','s3://eip-demo/datasets/images/v1','seed'),('00000000-0000-0000-0000-000000000403','00000000-0000-0000-0000-000000000001','dataset','demo-labels','v1','sha256:demo-labels-v1','sha256:demo-labels-v1','s3://eip-demo/datasets/labels/v1','seed') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO training_jobs (id,project_id,name,template_version,image_digest,code_version,dataset_version,resource_request,queue_id,state,output_uri) VALUES ('00000000-0000-0000-0000-000000000501','00000000-0000-0000-0000-000000000001','resnet50-demo','v1','sha256:demo-image-v1','git:demo-main','demo-images:v1','{"cpu":8,"memoryMb":32768,"gpu":1}','00000000-0000-0000-0000-000000000301','succeeded','s3://eip-demo/runs/resnet50-demo'),('00000000-0000-0000-0000-000000000502','00000000-0000-0000-0000-000000000001','yolov8-demo','v1','sha256:demo-image-v1','git:demo-main','demo-images:v1','{"cpu":8,"memoryMb":32768,"gpu":1}','00000000-0000-0000-0000-000000000302','running','s3://eip-demo/runs/yolov8-demo'),('00000000-0000-0000-0000-000000000503','00000000-0000-0000-0000-000000000001','pending-demo','v1','sha256:demo-image-v1','git:demo-main','demo-images:v1','{"cpu":4,"memoryMb":16384,"gpu":1}','00000000-0000-0000-0000-000000000303','pending_validation','s3://eip-demo/runs/pending-demo') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO training_job_assets (job_id,asset_version_id) VALUES ('00000000-0000-0000-0000-000000000501','00000000-0000-0000-0000-000000000401'),('00000000-0000-0000-0000-000000000501','00000000-0000-0000-0000-000000000402'),('00000000-0000-0000-0000-000000000502','00000000-0000-0000-0000-000000000401'),('00000000-0000-0000-0000-000000000502','00000000-0000-0000-0000-000000000402'),('00000000-0000-0000-0000-000000000503','00000000-0000-0000-0000-000000000401'),('00000000-0000-0000-0000-000000000503','00000000-0000-0000-0000-000000000403') ON CONFLICT DO NOTHING`,
		`INSERT INTO experiment_runs (id,training_job_id,parameters,metrics,logs_uri,artifacts_uri,best) VALUES ('00000000-0000-0000-0000-000000000601','00000000-0000-0000-0000-000000000501','{"epochs":20,"batchSize":64}','{"accuracy":0.94,"loss":0.08}','s3://eip-demo/runs/resnet50-demo/logs','s3://eip-demo/runs/resnet50-demo/artifacts',true) ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO model_versions (id,experiment_run_id,version,format,signature,metrics_summary,approval_status,lifecycle) VALUES ('00000000-0000-0000-0000-000000000701','00000000-0000-0000-0000-000000000601','v1','onnx','sha256:demo-model-v1','{"accuracy":0.94}','approved','active') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO online_services (id,project_id,model_version_id,name,image_digest,resource_spec,endpoint,status,version_policy) VALUES ('00000000-0000-0000-0000-000000000801','00000000-0000-0000-0000-000000000001','00000000-0000-0000-0000-000000000701','resnet50-online','sha256:demo-image-v1','{"cpu":4,"memoryMb":8192,"gpu":1}','https://resnet50.demo.eip.local','running','fixed') ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO audit_events (id,actor,project_id,resource_type,resource_id,action,after,request_id) VALUES ('00000000-0000-0000-0000-000000000901','seed','00000000-0000-0000-0000-000000000001','project','00000000-0000-0000-0000-000000000001','seeded','{"status":"active"}','seed-001') ON CONFLICT (id) DO NOTHING`,
	}
	for _, statement := range statements {
		if _, err := tx.Exec(statement); err != nil { return fmt.Errorf("seed: %w", err) }
	}
	return tx.Commit()
}
