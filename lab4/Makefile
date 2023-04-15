.PHONY: all server client-single client-multi clean

all: server client-single client-multi

server:
	go build -o ./server/server ./server/main.go

client-single:
	go build -o ./client-single/client ./client-single/main.go

client-multi:
	go build -o ./client-multi/client ./client-multi/main.go

clean:
	rm -f ./server/server ./client-single/client ./client-multi/client
