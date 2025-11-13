# usersvc (Hexagonal Architecture)

A minimal Go microservice that demonstrates **hexagonal architecture** (ports & adapters):

- **Driving adapter**: gRPC (`CreateUser`)
- **Application**: `CreateUser` use case (pure business workflow)
- **Driven adapters**: mocked repository (in-memory) and mocked publisher (stdout)
- **Kubernetes-ready** manifests
- **buf/protoc**-based codegen

## Hex Mapping
```
cmd/usersvc                # composition root (wires the hex)
internal/
  domain/                  # entities, value objects, domain errors
    user/
      entity.go
  app/                     # use cases + ports (interfaces)
    ports.go
    user/
      create_user.go
  adapters/
    driving/grpc/          # gRPC server -> calls use case
      server.go
    driven/
      repo/memory/         # implements UserRepository
        repo.go
      pub/stdout/          # implements UserPublisher
        publisher.go
api/
  proto/user/v1/user.proto
  gen/go/...               # generated (make proto)
k8s/deployment.yaml
```

## API
```
rpc CreateUser(CreateUserRequest) returns (CreateUserResponse)
```

## Run locally
```bash
make proto
make tidy
make test
make run
```

Invoke:
```bash
grpcurl -plaintext -d '{"name":"Alice","email":"alice@example.com"}'   localhost:8080 user.v1.UserService/CreateUser
```

## Swap Adapters (Production)
- Replace `repo/memory` with MongoDB adapter implementing `app.UserRepository`.
- Replace `pub/stdout` with Kafka adapter implementing `app.UserPublisher`.
- No changes to domain or use case code required.
