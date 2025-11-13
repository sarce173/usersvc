FROM golang:1.22 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/usersvc ./cmd/usersvc

FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/usersvc /usr/local/bin/usersvc
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/usersvc"]
