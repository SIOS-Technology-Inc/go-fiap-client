FROM golang:1.21.7

RUN apt update \
  && apt install -y vim

ENV GO111MODULE on
WORKDIR /go/src/work

# install go tools
RUN go install golang.org/x/tools/gopls@v0.14.2 &&\
  go install github.com/cweill/gotests/...@v1.6.0 &&\
  go install github.com/fatih/gomodifytags@v1.16.0 &&\
  go install github.com/josharian/impl@v1.3.0 &&\
  go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 &&\
  go install github.com/go-delve/delve/cmd/dlv@v1.22.0 &&\
  go install honnef.co/go/tools/cmd/staticcheck@v0.4.5 &&\
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.0
