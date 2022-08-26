FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /go/src/yujun/

COPY chatwork-service/. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/yujun/ .

CMD ["./prometheus-chatwork"]
