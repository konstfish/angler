FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main main/main.go

FROM alpine:latest

ENV GIN_MODE=release

RUN apk --no-cache add ca-certificates wget curl

WORKDIR /

RUN mkdir -p db

COPY container-init.sh .
COPY --from=builder /app/main .

CMD ["./container-init.sh"]