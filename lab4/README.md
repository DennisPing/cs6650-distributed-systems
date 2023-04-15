# Lab 4

Dennis Ping

## How to Build
```
make
```

## How to Run

### Server
```
./server/server 
```

### Single-threaded Client
```
./client-single/client -h localhost -p 12031
```

### Multi-threaded Client
```
./client-multi/client -h localhost -p 12031
```

## Multi-threaded Server

- [x] Worker pool for max concurrent connections = 20
- [x] Simulated hard work of 100ms

The client is able to make up to 100 initial connection attempts, however the server's worker pool size caps out at 20. The 20 worker threads will cycle through incoming requests as fast as possible.

Each worker thread sleeps for 100ms to simulate some server-side work.

## Single-threaded Client

Connect to the server, send the clientID, and wait for the server to send back the response.
Then close the connection.

## Multi-threaded Client

- [x] Number of threads = 1000
- [x] Max concurrent connections = 100
- [x] Connection retry and exponential backoff

The assignment asked us to use a cyclical barrier to synchronize all the threads, but Go's waitgroup is also suitable here. By using a waitgroup, the main program will wait for all goroutines to complete their tasks before terminating. Since each goroutine does not depend the result of the others, each one can disconnect from the server as soon as it finishes its task, thus releasing resources on the server side.

## UDP Multi-threaded Server

TODO

## UDP Multi-threaded Client

TODO

Will reuse a pool of connections rather than making & closing each time

![img](https://miro.medium.com/v2/resize:fit:1400/format:webp/1*xe4DmSW7U1PNY8vzryKZ6Q.png)