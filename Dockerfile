FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /fetch-bookmeter .

FROM alpine:3.19

COPY --from=builder /fetch-bookmeter /fetch-bookmeter

ENTRYPOINT ["/fetch-bookmeter"]
