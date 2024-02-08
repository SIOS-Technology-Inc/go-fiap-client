FROM golang:1.21

RUN apt update \
  && apt install -y vim # 不要な場合は削除してください

ENV GO111MODULE on
WORKDIR /go/src/work

# install go tools（自動補完等に必要なツールをコンテナにインストール）
RUN go install golang.org/x/tools/gopls@latest &&\
  go install github.com/cweill/gotests/...@latest &&\
  go install github.com/fatih/gomodifytags@latest &&\
  go install github.com/josharian/impl@latest &&\
  go install github.com/haya14busa/goplay/cmd/goplay@latest &&\
  go install github.com/go-delve/delve/cmd/dlv@latest &&\
  go install honnef.co/go/tools/cmd/staticcheck@latest &&\
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
