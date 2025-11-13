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

