# ------------------------------------------------------------
# Build stage
# ------------------------------------------------------------

FROM golang:1.25.3 AS builder

WORKDIR /app

# Copy dependency files first so Docker can cache dependencies.
COPY go.mod go.sum ./

RUN go mod download

# Copy the complete backend source.
COPY . .

# Build a Linux binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /rushd-server ./cmd/server


# ------------------------------------------------------------
# Runtime stage
# ------------------------------------------------------------

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /rushd-server /rushd-server

EXPOSE 8080

ENTRYPOINT ["/rushd-server"]