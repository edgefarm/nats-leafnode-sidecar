FROM golang:1.17 AS builder
WORKDIR /code
COPY . .
RUN make registry

FROM alpine:3.14
COPY --from=builder /code/bin/registry /registry
ENTRYPOINT ["/registry"]
