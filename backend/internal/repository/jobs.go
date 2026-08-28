package repository

import (
	"errors"
	"sync"
	"eip-platform/backend/internal/domain"
)

var jobStore = struct{sync.Mutex; jobs map[string]domain.TrainingJob; idem map[string]string}{jobs:map[string]domain.TrainingJob{}, idem:map[string]string{}}

func ListJobs(projectID string) []domain.TrainingJob { jobStore.Lock(); defer jobStore.Unlock(); out:=[]domain.TrainingJob{}; for _,j:=range jobStore.jobs {if j.ProjectID==projectID {out=append(out,j)}}; return out }
func CreateJob(job domain.TrainingJob, key string) (domain.TrainingJob,error) { jobStore.Lock(); defer jobStore.Unlock(); if key!="" {if id:=jobStore.idem[key];id!="" {return jobStore.jobs[id],nil}}; if job.ID==""||job.ProjectID==""||job.ImageDigest==""||job.CodeVersion==""||job.DatasetVersion=="" {return domain.TrainingJob{},errors.New("immutable asset references required")}; job.State=domain.JobPendingValidation; jobStore.jobs[job.ID]=job;if key!="" {jobStore.idem[key]=job.ID};return job,nil }
func UpdateJobState(id string, state domain.JobState) (domain.TrainingJob,error) { jobStore.Lock(); defer jobStore.Unlock(); j,ok:=jobStore.jobs[id];if !ok{return j,errors.New("job not found")};j.State=state;jobStore.jobs[id]=j;return j,nil }
func GetJob(id string)(domain.TrainingJob,bool){jobStore.Lock();defer jobStore.Unlock();j,ok:=jobStore.jobs[id];return j,ok}
