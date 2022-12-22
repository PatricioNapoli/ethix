FROM golang:1.18.2-alpine as builder

RUN apk add --no-cache make gcc

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

WORKDIR /go/src/ethix

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.mod
RUN go mod download

COPY . .

ENTRYPOINT ["/tini", "-g", "--"]
CMD ["./scripts/test.sh"]
