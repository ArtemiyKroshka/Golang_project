# Step 1: Build the Go application
FROM golang:1.22.3 AS builder

WORKDIR /go_project/

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_project/go_project main.go

# Step 2: Use a minimal image to run the Go application
FROM alpine:latest

ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=11111
ENV DB_NAME=golang_project_db
ENV DB_HOST=db

# Copy the built Go binary
COPY --from=builder /go_project/go_project /usr/local/bin/go_project

# Copy the templates directory
COPY --from=builder /go_project/internal/templates /internal/templates

# Set the entrypoint for the Go application
CMD ["/usr/local/bin/go_project"]
