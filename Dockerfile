FROM golang:1.17 as gobuild
COPY . nginx-watcher
WORKDIR /go/nginx-watcher/
ENV CGO_ENABLED=0
RUN go build -o watcher main.go


FROM scratch
COPY --from=gobuild /go/nginx-watcher/watcher /watcher
ENTRYPOINT ["/watcher"]