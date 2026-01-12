package job

import "time"

// JobPayload is a flexible container for job-specific data
type JobPayload map[string]any
type JobPriority int

const (
	Low JobPriority = iota + 1
	Medium
	High
)

// Job represents a unit of work in the worker queue
type Job struct {
	Id              string      // unique identifier for the job
	Type            string      // job type (e.g., send_email, generate_report)
	Payload         JobPayload  // data required to execute the job
	State           JobState    // current state of the job
	Attempts        int         // number of execution attempts so far
	MaxAttempts     int         // maximum number of retries allowed
	RunAt           *time.Time  // scheduled execution time
	LastRunAt       *time.Time  // timestamp of the last execution attempt
	FinishedAt      *time.Time  // timestamp of successful completion
	ErrorMessage    *string     // reason for last failure (if any)
	Priority        JobPriority // priority for ordering jobs
	IndempotencyKey string      // ensures safe retries without duplication
	CreatedAt       time.Time   // creation timestamp
	UpdatedAt       *time.Time  // last update timestamp
}
