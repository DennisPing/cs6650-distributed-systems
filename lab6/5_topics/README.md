# Subscribing to certain topics

Consumers can subscribe to any number of topics

## Concept

Route based on the location and the log severity, depending on what the consumer wants to listen for.

## Example

### Producer
```
go run main.go
>> boston.info Some normal message
>> boston.error Boston housing is too damn expensive
>> seattle.warn The weather is rainy here
>> nyc.error The datacenter in NYC crashed!
>> nyc.warn The traffic is bad today
```

### Start consumer 1 (info subscriber)
```
go run main.go "*.info"
2023/04/30 22:43:12 Received: Some normal message
```

### Start consumer 2 (east coast subscriber)
```
go run main.go "boston.*" "nyc.*"
2023/04/30 22:43:12 Received: Some normal message
2023/04/30 22:43:17 Received: Boston housing is too damn expensive
2023/04/30 22:43:25 Received: The datacenter in NYC crashed!
2023/04/30 22:43:29 Received: The traffic is bad today
```

### Start consumer 3 (warn and error subscriber)
```
go run main.go "*.warn" "*.error"
2023/04/30 22:43:17 Received: Boston housing is too damn expensive
2023/04/30 22:43:22 Received: The weather is rainy here
2023/04/30 22:43:25 Received: The datacenter in NYC crashed!
2023/04/30 22:43:29 Received: The traffic is bad today
```