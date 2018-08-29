FROM golang:1.10
COPY . /go/src/github.com/AliyunContainerService/open-service-broker-alibabacloud
WORKDIR /go/src/github.com/AliyunContainerService/open-service-broker-alibabacloud
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /alibabacloud-servicebroker .

FROM alpine:3.6
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/' /etc/apk/repositories
RUN apk update && apk add ca-certificates
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=0 /alibabacloud-servicebroker /alibabacloud-servicebroker
CMD ["/alibabacloud-servicebroker", "-logtostderr"]
