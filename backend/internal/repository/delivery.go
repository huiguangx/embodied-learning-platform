package repository

import("errors";"sync";"eip-platform/backend/internal/domain")
var deliveryStore=struct{sync.Mutex;runs map[string]domain.ExperimentRun;models map[string]domain.ModelVersion;services map[string]domain.OnlineService;audits []domain.AuditEvent}{runs:map[string]domain.ExperimentRun{},models:map[string]domain.ModelVersion{},services:map[string]domain.OnlineService{}}
func EnsureExperiment(j domain.TrainingJob) domain.ExperimentRun{deliveryStore.Lock();defer deliveryStore.Unlock();for _,r:=range deliveryStore.runs{if r.TrainingJobID==j.ID{return r}};r:=domain.ExperimentRun{ID:"run-"+j.ID,TrainingJobID:j.ID,Metrics:[]byte(`{"accuracy":0.94}`),ArtifactsURI:j.OutputURI,Best:true};deliveryStore.runs[r.ID]=r;return r}
func Runs(projectID string)[]domain.ExperimentRun{deliveryStore.Lock();defer deliveryStore.Unlock();out:=[]domain.ExperimentRun{};for _,r:=range deliveryStore.runs{if j,ok:=GetJob(r.TrainingJobID);ok&&j.ProjectID==projectID{out=append(out,r)}};return out}
func RegisterModel(m domain.ModelVersion)(domain.ModelVersion,error){deliveryStore.Lock();defer deliveryStore.Unlock();if _,ok:=deliveryStore.runs[m.ExperimentRunID];!ok{return m,errors.New("experiment run required")};m.ID="model-"+m.ExperimentRunID;deliveryStore.models[m.ID]=m;return m,nil}
func Models()[]domain.ModelVersion{deliveryStore.Lock();defer deliveryStore.Unlock();out:=[]domain.ModelVersion{};for _,m:=range deliveryStore.models{out=append(out,m)};return out}
func CreateService(s domain.OnlineService)(domain.OnlineService,error){deliveryStore.Lock();defer deliveryStore.Unlock();m,ok:=deliveryStore.models[s.ModelVersionID];if !ok||m.ApprovalStatus!="approved"{return s,errors.New("published model required")};s.ID="svc-"+s.Name;s.Status="running";deliveryStore.services[s.ID]=s;return s,nil}
func Services(projectID string)[]domain.OnlineService{deliveryStore.Lock();defer deliveryStore.Unlock();out:=[]domain.OnlineService{};for _,s:=range deliveryStore.services{if s.ProjectID==projectID{out=append(out,s)}};return out}
func Audit(e domain.AuditEvent){deliveryStore.Lock();defer deliveryStore.Unlock();deliveryStore.audits=append(deliveryStore.audits,e)}
