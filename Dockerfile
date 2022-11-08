FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY smib-vault-client/smib-vault-client .
RUN env GOOS=linux GOARCH=amd64 go build -o smib-vault-client
CMD ["/app/smib-vault-client"]

############################
FROM scratch
WORKDIR /app
COPY --from=builder /app/smib-vault-client /app/smib-vault-client
CMD ["/app/smib-vault-client"]
ENTRYPOINT ["/app/smib-vault-client"]