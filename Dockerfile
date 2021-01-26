FROM golang:1.15.7 as builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -o /go/bin/main -ldflags '-s -w'

FROM alpine as runner
#FROM golang:1.15.7 as runner

ENV SLACK_TOKEN="" \
    RUTTY_API_URL="localhost:3000/"

COPY --from=builder /go/bin/main /app/main
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/app/main" ]