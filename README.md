# Pulse

## Project Status: Ongoing Development

Pulse is a **durable worker queue and task scheduler** for running background jobs in real-world backend systems. The project is under active development and focuses on correctness, reliability, and observability.

Pulse is designed to handle asynchronous workloads that must survive crashes, support retries, and provide full operational visibility.

---

## Key Features

* **Database-backed job queue** for durability and visibility
* **Leader-elected scheduler** for priority and type-based job selection
* **Stateless workers** that claim and execute jobs atomically
* **Explicit job lifecycle** with retries, failure handling, and idempotency
* **Observability and inspection** for operational safety
* Designed for **correctness-first workloads**

---

## Getting Started (Minimal)

1. Start your PostgreSQL database.
2. Run the scheduler:

```bash
go run cmd/scheduler/main.go
```

3. Run one or more workers:

```bash
go run cmd/worker/main.go
```

4. Submit jobs via the API:

```bash
curl -X POST http://localhost:8080/jobs -d '{"type":"email","priority":1}'
```

*Note: These are placeholder commands. The system is under active development.*

---

## Design Philosophy & Trade-Offs

Pulse is **correctness-first, observable, and recoverable**. Some key design trade-offs include:

* **Database-backed queue vs broker**: ensures durability and easier debugging, but limits peak throughput.
* **Single leader scheduling**: simplifies prioritization and fairness, but can become a bottleneck under very high load.
* **Workers pull jobs rather than leader assigning them**: improves fault tolerance and recovery, avoids tight coupling.
* **Predictability over cleverness**: favors boring, maintainable design over complex optimizations.

For full architecture, detailed trade-offs, and failure-handling scenarios, see [docs/DESIGN.md](docs/DESIGN.md).

---

## Intended Usage

Pulse is designed to be embedded into larger systems where it powers:

* Notifications and messaging
* Scheduled tasks and cron replacements
* Retrying failed operations
* Asynchronous workflows

It is **not a standalone application** and deliberately avoids workflow DSLs or generic messaging semantics.

---

## Contributing

Contributions are welcome, but please focus on:

* Correctness and reliability
* Observability and operational safety
* Maintaining explicit job lifecycle semantics

---

## Disclaimer

Pulse is an **ongoing experimental project**. It is not a replacement for mature message brokers or workflow engines. It is intended for learning, experimentation, and controlled production usage.
