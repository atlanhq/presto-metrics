FROM golang:1.13.4 as builder

MAINTAINER arpit@atlan.com
COPY . "/go/src/github.com/atlanhq/presto-metrics/"

WORKDIR "/go/src/github.com/atlanhq/presto-metrics"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /presto_metrics

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /presto_metrics .
COPY --from=builder /go/src/github.com/atlanhq/presto-metrics/entrypoint.sh .
RUN  chmod +x entrypoint.sh
EXPOSE 9483
ENTRYPOINT ["./entrypoint.sh"]
