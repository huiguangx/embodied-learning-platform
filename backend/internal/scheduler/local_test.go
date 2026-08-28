package scheduler
import("context";"testing";"eip-platform/backend/internal/domain";"eip-platform/backend/internal/repository")
func TestLocalSchedulerCompletesJob(t *testing.T){j,_:=repository.CreateJob(domain.TrainingJob{ID:"sched-test",ProjectID:"p",ImageDigest:"i",CodeVersion:"c",DatasetVersion:"d"},"");if err:=(Local{}).Submit(context.Background(),j);err!=nil{t.Fatal(err)};got,_:=repository.GetJob(j.ID);if got.State!=domain.JobSucceeded{t.Fatalf("got %s",got.State)}}
