# Build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . ./

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g main.go -d cmd/api --parseDependency --parseInternal -o docs

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12
COPY --from=builder /bin/api /api
ENV API_PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]