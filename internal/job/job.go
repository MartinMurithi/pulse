package job

import (
	"fmt"
	"time"
)

type JobPayload struct {
	JobDescription string
}

type Job struct {
	ID      int
	Type    string
	Payload JobPayload
}

// Simulate a Job Producer

func JobProducer() {

	for i := 0; i <= 10; i++ {
		job := Job{
			ID:   i,
			Type: "test_job",
			Payload: JobPayload{
				JobDescription: fmt.Sprintf("job description with id : %d\n", i),
			},
		}

		fmt.Printf("submitting job [%d] to queue : job desc %s\n", job.ID, job.Type)

		// Simulate Delay to show incoming job rate
		time.Sleep(2 * time.Second)
	}

}

// Store the Generated Jobs in a queue

