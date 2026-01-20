# Pulse: Detailed Design Document

## Project Status: Ongoing Development

Pulse is a **durable worker queue and task scheduler** designed for running background jobs in real-world backend systems. The project is under active development and focuses on correctness, reliability, and observability. Pulse is intended for asynchronous workloads that must survive crashes, support retries, and provide full operational visibility.

---

## 1. Problem Statement

Most systems eventually require background execution for tasks such as:

* Sending emails or notifications
* Retrying failed external API calls
* Running delayed or scheduled jobs
* Processing webhooks
* Generating reports or exports

Ad-hoc solutions often fail due to:

* Lost jobs on process restart
* Duplicate execution without idempotency
* Unbounded goroutines causing resource exhaustion
* Misuse of cron jobs as retry mechanisms
* No visibility into failures or retries

Pulse enforces **durability, explicit state transitions, and controlled execution** to address these problems.

---

## 2. Core Design Principles

* **Durability first**: No job exists only in memory.
* **Explicit state transitions**: Jobs always occupy a single well-defined state.
* **Separation of concerns**: Scheduling and execution are independent.
* **Failure is expected**: Retries and backoff are first-class concepts.
* **Predictable and boring by design**: Correctness and observability are prioritized over cleverness.

---

## 3. High-Level Architecture

Pulse consists of four main components:

1. **Producer (service layer)**: Creates jobs as part of normal application workflows.
2. **Persistent job store**: Stores jobs and state transitions durably using PostgreSQL.
3. **Scheduler**: Determines when a job becomes runnable, promoting it to a `pending` state.
4. **Workers**: Execute jobs and record outcomes. Workers claim jobs atomically to ensure exclusive execution.

Each component has a single responsibility and minimal coupling with others.

---

## 4. Job Model and Lifecycle

Jobs are modeled as a **state machine**, enabling retries, inspection, and operational safety.

### Job States

* `scheduled` – waiting until `run_at`
* `pending` – eligible for execution
* `running` – currently being processed
* `retrying` – failed, waiting for backoff
* `completed` – finished successfully
* `failed` – terminal failure
* `dead` – exceeded retry limit (DLQ)

### Execution Flow

1. **Job Creation (Producer Phase)**

   * Job is validated and persisted in the database
   * Includes metadata: `job_id`, `job_type`, `payload`, `state`, `run_at`, `attempt`, `max_attempts`, `priority`, `idempotency_key`

2. **Scheduling**

   * Scheduler selects `scheduled` or `retrying` jobs whose `run_at` <= now
   * Promotes eligible jobs to `pending`

3. **Dispatching**

   * Workers poll for `pending` jobs
   * Claim jobs using database-level locking
   * Transition state to `running`

4. **Execution**

   * Worker dispatches job to the corresponding handler
   * Execution may be in-process or delegated via HTTP/gRPC

5. **Failure Handling and Retries**

   * On failure, backoff is calculated
   * Job moves to `retrying`, `run_at` updated
   * If retries exhausted, job moves to `dead` for inspection or manual requeue

6. **Graceful Shutdown**

   * Workers stop pulling new jobs, finish or safely release current jobs, and respect context cancellation

7. **Idempotency**

   * Jobs include `idempotency_key`; handlers must be safe to execute multiple times

---

## 5. Leader-Elected Scheduling

Leader responsibilities:

* Determine which jobs are eligible to run
* Promote jobs from `scheduled` or `retrying` to `pending`
* Enforce prioritization and fairness

Workers claim jobs independently; the leader does **not** assign jobs to specific workers.

This approach simplifies correctness and recovery while avoiding a single bottleneck in execution.

---

## 6. Observability

Pulse enables inspection at all times:

* Query running, pending, and failed jobs
* Track retries and backoff
* Capture execution duration and outcomes
* Optional admin or CLI interfaces for operational visibility

---

## 7. Technology Choices

* **Go**: Explicit concurrency, predictable performance, clarity in failure handling
* **PostgreSQL**: Strong consistency, row-level locking, transactional state transitions, easy inspection

---

## 8. Design Constraints

Pulse intentionally avoids:

* Job logic inside HTTP handlers
* `time.Sleep` for scheduling
* Goroutine-per-request execution
* Scheduler logic mixed into repositories

---

## 9. Intended Usage

Pulse is designed to be embedded into larger systems for:

* Notifications
* Retries
* Scheduled tasks
* Asynchronous workflows

It is **not a standalone application** and is still under active development.

---

## 10. Scope and Limitations

**Priorities:** Correctness, durability, observability, operational safety

**Deliberate omissions:** Over-generalization, messaging semantics, workflow DSLs

**Known limitations:**

* Single-leader scheduling can be a bottleneck at scale
* Database-backed queue limits maximum throughput
* Operational discipline required for production use

---

## 11. Future Enhancements

* Job partitioning by type or tenant
* Multiple leaders per partition
* Pluggable queue backends
* Rate limiting and SLA enforcement
