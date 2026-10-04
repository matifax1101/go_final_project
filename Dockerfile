FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /src/scheduler /app/scheduler
COPY --from=builder /src/web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

CMD ["/app/scheduler"]
