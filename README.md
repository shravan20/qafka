# Qafka

## Overview
A self-service, developer-friendly REST API layer for Apache Kafka to simplify queue management, message handling, and extended features like DLQ and retry queues.

This open-source tool enables product teams, platform teams, and microservices to create and manage Kafka queues without needing to understand Kafka internals.

---

## 🌟 Key Features

### Core Queue Management APIs
- `POST /queues` – Create a topic with config
- `GET /queues/:name` – Fetch topic metadata
- `DELETE /queues/:name` – Delete topic
- `PATCH /queues/:name/config` – Update topic config (retention, partitions, etc.)
- `GET /queues` – List all topics (with optional filters like tenant, tags, etc.)

### Message Operations
- `POST /queues/:name/messages` – Produce message to a topic
- `GET /queues/:name/messages` – Read messages (seek by offset/timestamp)
- `POST /queues/:name/seek` – Rewind consumer to offset or timestamp

### Advanced Message Handling
- **Dead Letter Queue (DLQ)**
  - Auto-create `<topic>-dlq` with separate retention
  - Route failed messages on retry exhaustion
- **Retry Queues**
  - Configure retry strategies (linear, exponential, fixed interval)
- **Message Expiry**
  - Optional TTL per message
- **Schema Validation**
  - Optional Avro/JSON schema enforcement

### Security & Governance
- Multi-tenant support with namespace isolation
- JWT/API Key-based auth with RBAC
- Access control for operations (create, produce, consume)
- Audit logs of queue events
- Quota enforcement (throughput, topic limit, etc.)

### Monitoring & Observability
- Topic lag monitoring per consumer group
- API for queue health status
- View last N messages via REST or UI
- Export Prometheus metrics
- Alert integrations (Slack, email, etc.)

### Dev Experience
- OpenAPI spec for REST
- SDKs in Node.js, Python, Go
- CLI tool: `kqctl`
- Web dashboard to manage and inspect queues
- Dry-run support for queue config validation

### Integrations
- GitOps-style declarative queue management (via config sync)
- Webhook triggers on DLQ / retries / threshold breaches
- Schema Registry support (Confluent/Apicurio)
- Pluggable DLQ backends (Kafka, S3, PostgreSQL)

### Future Roadmap Ideas
- Multi-cluster Kafka routing
- SLA Policies on queues (e.g., delivery deadlines)
- Topic versioning
- Message reprocessing UI with filters

---

## 🔧 Tech Stack Suggestions

- **Backend:** Node.js (Express/Fastify) / Go / Python (FastAPI)
- **Kafka Client:** kafkajs / spring-kafka / confluent-kafka-python
- **DB (for metadata):** PostgreSQL / SQLite (for dev)
- **Auth:** JWT / OAuth2 / API keys
- **Monitoring:** Prometheus + Grafana
- **UI (Optional):** React + Tailwind

---

## 💡 Example Use Cases

- Internal tooling for platform/infra teams
- CI/CD pipelines that need ephemeral queues
- Microservice-driven apps with custom retry/DLQ logic
- External apps needing simplified Kafka interaction

---

## 📘 License

Open Source – Apache 2.0
