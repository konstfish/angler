FROM golang:1.21-alpine as builder

WORKDIR /app

COPY geoip-api/go.mod geoip-api/go.sum ./
RUN go mod download

COPY go.work .
COPY geoip-api geoip-api

COPY shared shared

COPY geoip-api geoip-api
COPY ingress ingress
COPY auth auth
COPY backend backend

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main geoip-api/main/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget curl

WORKDIR /

RUN mkdir -p db

COPY --from=builder /app/main .
COPY container-init.sh .

CMD ["./container-init.sh"]