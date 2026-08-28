package scheduler

import("context";"eip-platform/backend/internal/domain";"eip-platform/backend/internal/repository")
type Local struct{}
func(Local)Submit(ctx context.Context,j domain.TrainingJob)error{for _,e:=range []string{"queue","allocate","start","succeed"}{s,err:=domain.Transition(j.State,e);if err!=nil{return err};j.State=s;if _,err=repository.UpdateJobState(j.ID,s);err!=nil{return err}};repository.EnsureExperiment(j);return nil}
func(Local)Cancel(ctx context.Context,id string)error{j,ok:=repository.GetJob(id);if !ok{return context.Canceled};s,err:=domain.Transition(j.State,"cancel");if err!=nil{return err};_,err=repository.UpdateJobState(id,s);return err}
