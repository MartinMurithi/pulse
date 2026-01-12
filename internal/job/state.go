package job

import "time"

// JobState represents possible states of a job
type JobState string

const (
	StatePending   JobState = "pending"   //Ready to Run
	StateScheduled JobState = "scheduled" //Waiting for RunAt
	StateRunning   JobState = "running"
	StateRetrying  JobState = "retrying"
	StateFailed    JobState = "failed"
	StateCompleted JobState = "completed"
	StateDead      JobState = "dead" // goes to DLQ(Dead Letter Queue)
)

// CanRun, checks if a job is eligible to be picked by a worker
func (j *Job) CanRun() bool {
	switch j.State {
	case StatePending:
		return true
	case StateScheduled:
		return j.RunAt != nil && j.RunAt.Before(time.Now())
	case StateRetrying:
		return true
	default:
		return false
	}
}

// Running, marks a job as currently running
func (j *Job) Running() {
	j.State = StateRunning
	j.ErrorMessage = nil
	now := time.Now()
	j.LastRunAt = &now
}

// MarkCompleted, marks a job as completed
func (j *Job) MarkCompleted() {
	j.State = StateCompleted
	now := time.Now()
	j.FinishedAt = &now
}

// MarkFailed handles failed jobs, retries and dead jobs
func (j *Job) MarkFailed(errMsg string) {
	j.Attempts++
	j.State = StateFailed
	j.ErrorMessage = &errMsg
	now := time.Now()
	j.FinishedAt = &now

	if j.Attempts < j.MaxAttempts {
		j.State = StateRetrying
	} else {
		j.State = StateDead
	}

}

// MarkScheduled marks a job as scheduled for a future run
func (j *Job) MarkScheduled(runAt time.Time) {
	j.State = StateScheduled
	j.RunAt = &runAt
}

// MarkDead marks a job as permanently failed (DLQ)
func (j *Job) MarkDead(errMsg string) {
	j.State = StateDead
	j.ErrorMessage = &errMsg
	now := time.Now()
	j.FinishedAt = &now
}
