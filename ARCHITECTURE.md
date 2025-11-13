# Architecture

# Hexagonal Architecture Notes

## Core ideas
- Domain model is **framework-agnostic**.
- Use cases orchestrate domain behavior and depend on **ports** (interfaces).
- Adapters implement ports, isolating infra (gRPC, Kafka, MongoDB).

## Kubernetes
- Multi-replica Deployment (gRPC on :8080), ClusterIP Service.
- gRPC health for probes; prefer grpc-health-probe or TCP probe.
- Non-secret config via ConfigMap; credentials in Secret.
- Use HPA on CPU/RPS; graceful termination with preStop.
- Distroless image, run as non-root.

## Scaling
- **gRPC**: multiple pods behind Service/mesh; apply keepalives, deadlines, TLS/mTLS.
- **Kafka**: idempotent producer, partition by `user_id`, RF=3, ISR>=2, schema registry, DLQ.
- **MongoDB**: replica set or Atlas; unique index on `email`; majority writes; tuned read prefs.

## Reliability
- Event publish is best-effort in this demo. For stronger guarantees, implement **Outbox**:
  - Write user + outbox record in one tx.
  - Background publisher reads outbox and publishes to Kafka with retries.
  - Mark outbox record as sent on success.

## CI/CD
- Run `buf generate`, unit tests, and static analysis in CI.
- Build minimal container; scan image; progressive rollout (canary).


## How it runs in Kubernetes
- **Container**: Distroless image, non-root user, gRPC on :8080.
- **Deployment**: `apps/v1 Deployment` with ≥2 replicas; **ClusterIP Service** on port 8080 (HTTP/2).
- **Probes**: gRPC health (or TCP) for readiness/liveness; graceful termination with small `preStop` delay.
- **Config**: Environment variables; future: ConfigMap for non-secrets, Secret for credentials (Mongo/Kafka).
- **Autoscaling**: HPA on CPU and/or custom metrics (e.g., RPS, p95 latency) via Prometheus/Adapter.
- **Observability**: Structured logs, OpenTelemetry traces/metrics (RPC latency, error rate, publish latency).
- **Security**: mTLS via service mesh (Istio/Linkerd) or ingress TLS; least privilege, seccomp, read-only fs.
- **Rollouts**: RollingUpdate (maxSurge=1, maxUnavailable=0); canary via mesh/Argo Rollouts.

## Scaling in Production

### gRPC (stateless app tier)
- Horizontal scale by adding pods; clients use round-robin/lb from Service/mesh.
- Enforce **deadlines** and **timeouts**; enable **keepalives** to clean up dead connections.
- Tune `MaxConcurrentStreams`; protect server with sane limits.
- Backward compatible protobuf evolution; versioned pkg (`user.v1`, later `user.v2`).

### Kafka (event publishing)
- Use managed Kafka (Confluent/MSK) or dedicated cluster.
- Topic: `user.created`, keys by `user_id` for per-user ordering.
- **Partitions**: start small (e.g., 6–12), scale with throughput/consumers; **RF=3**, **minISR=2**.
- **Producer**: idempotent, acks=all, retries with exponential backoff, batching + compression (zstd).
- **Schema**: Avro/Protobuf + Schema Registry; backward-compatible evolution.
- **Reliability**: prefer **Outbox** pattern for atomic write+publish; DLQ for poison messages.
- **Observability**: delivery latency, batch sizes, error rates; alert on prolonged retry spikes.

### MongoDB (persistence, future adapter)
- Managed Atlas or self-hosted **replica set** (1 primary, 2 secondaries).
- **Indexes**: unique on `email`, plus createdAt; size working set so hot indexes fit RAM.
- **Write concern**: `w=majority`, `j=true`; balance latency vs durability.
- **Read prefs**: primary for writes; secondaries for read-heavy/reporting if needed.
- **Backups**: continuous backups and periodic restore drills; TTL indexes for ephemeral data if applicable.
- **Migrations**: tool-driven (e.g., goose, Atlas CLI); blue/green or rolling with backward compatibility.

## Hexagonal Architecture Mapping
- **Driving adapter**: `internal/adapters/driving/grpc` (handles transport-level concerns only).
- **Application**: `internal/app/user` (use case orchestrating domain + ports).
- **Domain**: `internal/domain/user` (entity + validation; no infra dependencies).
- **Driven adapters**: `repo/memory`, `pub/stdout` implement `UserRepository`/`UserPublisher` ports.
- Swap-in Mongo/Kafka adapters without touching domain/use case code.

## How I used AI (brief)
- **Scaffolding**: Generated project structure, gRPC proto, and Kubernetes/Docker boilerplate.
- **Code generation**: Wrote hexagonal layers (domain, ports, adapters) and unit tests.
- **Docs**: Produced interview-focused architecture notes and a quickstart.
- **Editing**: I reviewed and adjusted naming, comments, and safety defaults (non-root, distroless).
