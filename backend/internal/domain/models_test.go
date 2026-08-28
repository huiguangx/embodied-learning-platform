package domain

import "testing"

func TestJobStates(t *testing.T) {
	for _, state := range []JobState{JobDraft, JobPendingValidation, JobQueued, JobAllocating, JobRunning, JobSucceeded, JobFailed, JobCancelled, JobStopped, JobTimeout} {
		if err := state.Validate(); err != nil { t.Errorf("%q: %v", state, err) }
	}
	if JobState("unknown").Valid() { t.Fatal("unknown state is valid") }
}
