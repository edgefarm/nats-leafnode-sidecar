FROM golang:1.18 AS builder
WORKDIR /code
COPY . .
RUN make registry

FROM alpine:3.16
COPY --from=builder /code/bin/registry /registry
CMD ["/registry"]
