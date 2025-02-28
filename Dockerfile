# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
# No go.sum file in the repo, so don't try to copy it
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jenkinscalc .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/jenkinscalc .
EXPOSE 8090
ENTRYPOINT ["./jenkinscalc"]