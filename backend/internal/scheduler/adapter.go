package scheduler

import("context";"eip-platform/backend/internal/domain")
type Scheduler interface{Submit(context.Context,domain.TrainingJob)error;Cancel(context.Context,string)error}
