# Build the manager binary
FROM golang:1.21.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev libc6-compat build-base libc-dev

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

#RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/

# Build
RUN GOOS=linux GOARCH=amd64 go build -o manager ./cmd/main.go


FROM alpine

WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]
