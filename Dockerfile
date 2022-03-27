FROM golang:1.18-alpine3.15 as builder

COPY . /github.com/tumarov/feeddy
WORKDIR /github.com/tumarov/feeddy

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:3.15

WORKDIR /root/
COPY --from=0 /github.com/tumarov/feeddy/bin/bot .
COPY --from=0 /github.com/tumarov/feeddy/configs configs/

EXPOSE 80

CMD ["./bot"]