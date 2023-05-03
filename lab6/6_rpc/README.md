# Remote Procedure Call (RPC)

## Concept

Multiple clients make RPC calls to a server which does work and returns the result back. This is a blocking action.

## Example

### Client
```
go run main.go 42
2023/05/02 20:46:49 Got: 267914296
```

### Server
```
go run main.go
2023/05/02 20:46:49 Computed fib(42) = 267914296
2023/05/02 20:46:49 Sending response: 267914296
```