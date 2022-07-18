FROM golang:1.18 AS builder
WORKDIR /code
COPY . .
RUN make client

FROM alpine:3.16
COPY --from=builder /code/bin/client /client
CMD ["/client"]
