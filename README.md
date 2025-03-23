# PoC - Protobuf + RabbitMQ

## Goal(s)

the ultimate goal is to prove out using protobufs without gRPC. To do this, we build a protobuf message and send it through RabbitMQ.

- Compile a protobuf message
- Send it through RabbitMQ
- Read it from RabbitMQ
- Print the message

## Getting Started

Spin up resources:

```
docker compose up
```

Run it:

```
make run
```
