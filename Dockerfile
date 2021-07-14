FROM golang:1.15 as build
ENV  GOPROXY=https://goproxy.cn,direct
WORKDIR $GOPATH/mcp-over-xds
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o=/bin/mcp-over-xds ./cmd

FROM alpine
COPY --from=build /bin/mcp-over-xds /usr/local/bin/mcp-over-xds
ENTRYPOINT ["/usr/local/bin/mcp-over-xds", "server"]
