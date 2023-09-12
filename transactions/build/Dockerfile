FROM golang:1.21 as builder

WORKDIR /app

RUN go install github.com/bufbuild/buf/cmd/buf@v1.26.1

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./buf*.yaml ./
COPY ./transactions ./transactions

RUN buf generate

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o serve ./transactions/cmd/serve

FROM scratch as production

COPY --from=builder /app/serve /bin/serve

CMD ["/bin/serve"]