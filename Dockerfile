
FROM  golang:1.18 as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

ADD . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o minCDN .

FROM alpine:latest AS production

COPY --from=builder /build/minCDN .

ENTRYPOINT ["./minCDN"]

CMD ["-UseEtcd","-EtcdAddress","127.0.0.1:10086","-MinioEndpoint","127.0.0.1:9000","-MinioAccessID","minioadmin","-MinioSecret","minioadmin"]


