# FILEPATH:/Users/yangyang/go/src/github.com/github-wyy/m-scheduler-extender/Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o m-scheduler-extender .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/m-scheduler-extender .
EXPOSE 8010
ENTRYPOINT ["./m-scheduler-extender"]
