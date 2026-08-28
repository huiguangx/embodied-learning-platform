package domain

import "fmt"

func Transition(current JobState, event string) (JobState, error) {
	if !current.Valid() { return "", fmt.Errorf("invalid current state") }
	switch event {
	case "validate": if current==JobDraft { return JobPendingValidation,nil }
	case "queue": if current==JobPendingValidation { return JobQueued,nil }
	case "allocate": if current==JobQueued { return JobAllocating,nil }
	case "start": if current==JobAllocating { return JobRunning,nil }
	case "succeed": if current==JobRunning { return JobSucceeded,nil }
	case "fail": if current==JobRunning||current==JobAllocating { return JobFailed,nil }
	case "cancel": if current==JobDraft||current==JobPendingValidation||current==JobQueued { return JobCancelled,nil }
	case "stop": if current==JobRunning { return JobStopped,nil }
	case "timeout": if current==JobRunning||current==JobQueued { return JobTimeout,nil }
	}
	return "", fmt.Errorf("event %q is not allowed from %q", event, current)
}
