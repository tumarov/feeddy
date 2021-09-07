FROM golang:1.17-alpine3.13 as builder

COPY . /github.com/tumarov/feeddy
COPY .env /github.com/tumarov/feeddy/.env
WORKDIR /github.com/tumarov/feeddy

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /github.com/tumarov/feeddy/bin/bot .
COPY --from=0 /github.com/tumarov/feeddy/.env .
COPY --from=0 /github.com/tumarov/feeddy/configs configs/

EXPOSE 80

CMD ["./bot"]