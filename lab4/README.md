# Lab 4

Dennis Ping

## Getting Started

### Build all binaries
```
make
```

### Run Server
```
./tcp-server/server 
```

### Run Single-threaded Client
```
./tcp-client-single/client -h localhost -p 12031
```

### Run Multi-threaded Client
```
./tcp-client-multi/client -h localhost -p 12031
```

## Multi-threaded Server

- [x] Worker pool size = 20

The client is able to make up to 100 initial connection attempts, however the server's worker pool size caps out at 20. The 20 worker threads will cycle through incoming requests as fast as possible.

## Single-threaded Client

Connect to the server, send the clientID, and wait for the server to send back the response.
Then close the connection. Only 1 message sent.

## Multi-threaded Client

- [x] Number of threads = 10,000
- [x] Max concurrent connections = 100
- [x] Connection retry and exponential backoff

The assignment asked us to use a cyclical barrier to synchronize all the threads, but Go's waitgroup is also suitable here. By using a waitgroup, the main program will wait for all goroutines to complete their tasks before terminating. Since each goroutine does not depend the result of the others, each one can disconnect from the server as soon as it finishes its task, thus releasing resources on the server side.

A max concurrent connections limit is needed in order to prevent a DDoS attack on the server.

## UDP Multi-threaded Server

- [x] Worker pool size = 20

The client is able to make up to 100 initial connection attempts, however the server's worker pool size caps out at 20. The 20 worker threads will cycle through incoming requests as fast as possible.

## UDP Multi-threaded Client

- [x] Number of threads = 10,000
- [x] Reusable connection pool = 100

Since UDP is a connectionless protocol, it makes sense to share a small pool of UDP connections (sockets) among a large number of threads.

Without this connection pool, you can actually exhaust the number of file descriptors on your computer. On my Mac, the limit is 256  (`ulimit -n`).

## Results

| Server             | Client             | Threads | Time Taken |
| ------------------ | ------------------ | ------- | ---------- |
| TCP multi-threaded | TCP multi-threaded | 10,000  | ~800ms     |
| UDP multi-threaded | UDP multi-threaded | 10,000  | ~300ms     |

## Worker Pool Diagram

![img](https://miro.medium.com/v2/resize:fit:1400/format:webp/1*xe4DmSW7U1PNY8vzryKZ6Q.png)