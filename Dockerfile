FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o load-tester ./cmd/load-tester

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/load-tester .

ENTRYPOINT ["./load-tester"]

LABEL maintainer="tiago-kimura"
LABEL description="Load Tester CLI"
LABEL version="1.0.0"