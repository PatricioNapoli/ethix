FROM golang:1.18.2-alpine as builder

RUN apk add --no-cache make gcc

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

RUN update-ca-certificates

ENV USER=svc
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/src/ethix

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.mod
RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static-debian11

EXPOSE 8080

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/src/ethix ./ethix
COPY --from=builder /tini /tini

ENTRYPOINT ["/tini", "-g", "--"]
CMD ["./ethix/bin/ethix"]
