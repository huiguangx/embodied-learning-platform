package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type JobState string

const (
	JobDraft JobState = "draft"
	JobPendingValidation JobState = "pending_validation"
	JobQueued JobState = "queued"
	JobAllocating JobState = "allocating"
	JobRunning JobState = "running"
	JobSucceeded JobState = "succeeded"
	JobFailed JobState = "failed"
	JobCancelled JobState = "cancelled"
	JobStopped JobState = "stopped"
	JobTimeout JobState = "timeout"
)

func (s JobState) Valid() bool {
	switch s {
	case JobDraft, JobPendingValidation, JobQueued, JobAllocating, JobRunning,
		JobSucceeded, JobFailed, JobCancelled, JobStopped, JobTimeout:
		return true
	}
	return false
}

func (s JobState) Validate() error {
	if !s.Valid() {
		return fmt.Errorf("invalid job state %q", s)
	}
	return nil
}

type Project struct {
	ID string `json:"id"`
	Organization string `json:"organization"`
	Name string `json:"name"`
	CostCenter string `json:"costCenter"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AssetVersion struct {
	ID string `json:"id"`
	ProjectID string `json:"projectId"`
	AssetType string `json:"assetType"`
	Name string `json:"name"`
	Version string `json:"version"`
	Digest string `json:"digest"`
	Checksum string `json:"checksum"`
	URI string `json:"uri"`
	Permissions json.RawMessage `json:"permissions,omitempty"`
	Source string `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

type ResourceGroup struct {
	ID string `json:"id"`
	ProjectID string `json:"projectId"`
	Name string `json:"name"`
	ClusterIDs []string `json:"clusterIds"`
	NodeType string `json:"nodeType"`
	CPUQuota int `json:"cpuQuota"`
	MemoryQuotaMB int `json:"memoryQuotaMb"`
	GPUQuota int `json:"gpuQuota"`
	Concurrency int `json:"concurrency"`
	BudgetTag string `json:"budgetTag"`
}

type Queue struct {
	ID string `json:"id"`
	ResourceGroupID string `json:"resourceGroupId"`
	Name string `json:"name"`
	Priority int `json:"priority"`
	Weight int `json:"weight"`
	Concurrency int `json:"concurrency"`
	MaxDurationSeconds int `json:"maxDurationSeconds"`
	SchedulingPolicy string `json:"schedulingPolicy"`
}

type TrainingJob struct {
	ID string `json:"id"`
	ProjectID string `json:"projectId"`
	Name string `json:"name"`
	TemplateVersion string `json:"templateVersion"`
	ImageDigest string `json:"imageDigest"`
	CodeVersion string `json:"codeVersion"`
	DatasetVersion string `json:"datasetVersion"`
	ResourceRequest json.RawMessage `json:"resourceRequest"`
	QueueID string `json:"queueId"`
	State JobState `json:"state"`
	OutputURI string `json:"outputUri"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ExperimentRun struct {
	ID string `json:"id"`
	TrainingJobID string `json:"trainingJobId"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
	Metrics json.RawMessage `json:"metrics,omitempty"`
	LogsURI string `json:"logsUri"`
	ArtifactsURI string `json:"artifactsUri"`
	Best bool `json:"best"`
	CreatedAt time.Time `json:"createdAt"`
}

type ModelVersion struct {
	ID string `json:"id"`
	ExperimentRunID string `json:"experimentRunId"`
	Version string `json:"version"`
	Format string `json:"format"`
	Signature string `json:"signature"`
	MetricsSummary json.RawMessage `json:"metricsSummary,omitempty"`
	ApprovalStatus string `json:"approvalStatus"`
	Lifecycle string `json:"lifecycle"`
	CreatedAt time.Time `json:"createdAt"`
}

type OnlineService struct {
	ID string `json:"id"`
	ProjectID string `json:"projectId"`
	ModelVersionID string `json:"modelVersionId"`
	Name string `json:"name"`
	ImageDigest string `json:"imageDigest"`
	ResourceSpec json.RawMessage `json:"resourceSpec"`
	Endpoint string `json:"endpoint"`
	Status string `json:"status"`
	VersionPolicy string `json:"versionPolicy"`
}

type AuditEvent struct {
	ID string `json:"id"`
	Actor string `json:"actor"`
	ProjectID string `json:"projectId"`
	ResourceType string `json:"resourceType"`
	ResourceID string `json:"resourceId"`
	Action string `json:"action"`
	Before json.RawMessage `json:"before,omitempty"`
	After json.RawMessage `json:"after,omitempty"`
	RequestID string `json:"requestId"`
	CreatedAt time.Time `json:"createdAt"`
}
