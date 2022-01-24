# nats-leafnode-sidecar

This repository contains the source of two components that are needed to operate edgefarm.network on the edgefarm application runtimes.

- `registry` and
- `client`

## `registry`

The `registry` is deployed as a sidecar container for the nats server and is responsible for registering the edgefarm application runtimes. During registration it configures the nats server configuration to allow the local running nats server to connect to the remote nats server. For this purpose each application gets its own credentials that need to get registered with the nats server.

## `client`

The `client` is deployed for each edgefarm.network resource that allows edgefarm.applications to interact with the remote nats server. The `client` registers itself on startup with the `registry` and passes its credentials that were generated from edgefarm.network operator.

## Building

Run `make` to build the binaries. Docker images are built and released via github actions.
