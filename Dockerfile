FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM alpine:3.22

WORKDIR /app
RUN mkdir -p /data
COPY --from=builder /src/scheduler /app/scheduler
COPY --from=builder /src/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["/app/scheduler"]
