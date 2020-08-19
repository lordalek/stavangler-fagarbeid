# Generate proto
```protoc -I order/ order/order.proto --go_out=plugins=grpc:order ```