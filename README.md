# Pulse
---
## What Pulse Is

Pulse is a **durable worker queue and task scheduler** for running background jobs in real-world backend systems.

It is designed for situations where work must:

* Run outside the request/response path
* Survive crashes and restarts
* Execute at a specific time or after a delay
* Retry safely on failure
* Be observable and debuggable

Pulse focuses strictly on **job orchestration**, not business logic.

---

## The Problem Pulse Solves

Most systems eventually need background execution for tasks such as:

* Sending emails or notifications
* Retrying failed external API calls
* Running delayed or scheduled jobs
* Processing webhooks
* Generating reports or exports

Common failure modes in ad-hoc solutions include:

* Jobs lost on process restart
* Duplicate execution without idempotency
* Unbounded goroutines
* Cron jobs misused as retry mechanisms
* No visibility into failures or retries

Pulse addresses these problems by enforcing **durability, explicit state, and controlled execution**.

---

## Core Design Principles

Pulse is built around the following principles:

* **Durability first**
  No job exists only in memory.

* **Explicit state transitions**
  Every job is always in exactly one state.

* **Separation of concerns**
  Scheduling and execution are independent.

* **Failure is expected**
  Retries and backoff are first-class concepts.

* **Boring by design**
  Predictability is more important than cleverness.

---

## System Components

Pulse consists of four independent actors:

1. **Producer (service layer)**
   Creates jobs as part of normal application workflows.

2. **Persistent job store**
   Stores jobs and their state transitions durably.

3. **Scheduler**
   Decides *when* a job becomes runnable.

4. **Workers**
   Execute jobs and record outcomes.

Each component has a single responsibility and minimal knowledge of the others.

---

## Job Model

Jobs are modeled as a **state machine**, not as fire-and-forget tasks.

### Job states

* `scheduled` – waiting until `run_at`
* `pending` – eligible for execution
* `running` – currently being processed
* `retrying` – failed, waiting for backoff
* `completed` – finished successfully
* `failed` – terminal failure
* `dead` – exceeded retry limit (DLQ)

This explicit lifecycle enables retries, inspection, and operational safety.

---

## End-to-End Execution Flow

### 1. Job creation (producer phase)

A service decides background work is required:

* Validates business rules
* Performs transactional work
* Persists a job record

At minimum, a job includes:

* `job_id`
* `job_type`
* `payload`
* `state`
* `run_at`
* `attempt`
* `max_attempts`
* `priority`
* `idempotency_key`

Once persisted, the job is guaranteed not to be lost.

---

### 2. Scheduling (time decision)

The scheduler runs independently and periodically:

* Selects jobs in `scheduled` or `retrying`
* Checks `run_at <= now`
* Atomically promotes jobs to `pending`

The scheduler never executes business logic.

---

### 3. Dispatching (queue consumption)

Workers:

* Poll for `pending` jobs
* Claim jobs using database-level locking
* Transition jobs to `running`
* Increment attempt counters

Only one worker can own a job at a time.

---

### 4. Execution (business logic)

The worker:

* Dispatches by `job_type`
* Executes the corresponding handler
* Captures duration and outcome

Pulse does not contain domain logic. Execution may be:

* In-process
* Delegated to another service via HTTP or gRPC

---

### 5. Failure handling and retries

If execution fails:

* Backoff is calculated
* Job transitions to `retrying`
* `run_at` is updated

If retries are exhausted:

* Job moves to `dead`
* Available for inspection or manual requeue

---

### 6. Graceful shutdown

On shutdown signals, workers:

* Stop pulling new jobs
* Finish or safely release current jobs
* Respect context cancellation

This ensures zero job loss during deployments.

---

## Idempotency Guarantees

Pulse assumes:

> **Any job may execute more than once.**

To support this, jobs include an `idempotency_key`, and systems using Pulse must ensure handlers are safe to retry.

This is a deliberate design constraint, not an implementation detail.

---

## Observability and Inspection

Pulse is designed to be inspectable at all times.

You should be able to answer:

* What is currently running?
* What is scheduled next?
* What failed and why?
* How many retries occurred?

This is achieved through:

* Explicit job states
* Structured logging
* Queryable job storage
* Optional admin or CLI interfaces

---

## Technology Choices

### Go

Used for:

* Explicit concurrency
* Simple deployment
* Clear failure semantics
* Predictable performance

The code favors clarity over abstraction.

---

### PostgreSQL

Used as the primary job store.

Chosen for:

* Strong consistency
* Row-level locking
* Transactional state transitions
* Easy inspection and recovery

Pulse does not depend on in-memory queues for correctness.

---

## Design Constraints (Enforced)

Pulse intentionally disallows the following patterns:

* Job logic in HTTP handlers
* `time.Sleep` for scheduling
* Goroutine-per-request execution
* Scheduler logic mixed into repositories

Violating these rules introduces silent failure modes.

---

## Intended Usage

Pulse is not a standalone application.

It is designed to be **embedded into larger systems**, where it acts as the background execution engine powering:

* Notifications
* Retries
* Scheduled tasks
* Asynchronous workflows

It can be used across multiple repositories and services.

---

## Scope

Pulse prioritizes:

* Correctness
* Durability
* Observability
* Operational safety

It deliberately avoids:

* Over-generalization
* Messaging semantics
* Workflow DSLs

---

## Status

Pulse is under active development with a focus on a minimal, production-correct core.

The initial goal is a system that can be trusted to run real background work in real systems.

---
