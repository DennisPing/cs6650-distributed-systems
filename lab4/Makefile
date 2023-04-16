.PHONY: all tcp-server tcp-client-single tcp-client-multi udp-server udp-client clean

all: tcp-server tcp-client-single tcp-client-multi udp-server udp-client

tcp-server:
	go build -o ./tcp-server/server ./tcp-server/main.go

tcp-client-single:
	go build -o ./tcp-client-single/client ./tcp-client-single/main.go

tcp-client-multi:
	go build -o ./tcp-client-multi/client ./tcp-client-multi/main.go

udp-server:
	go build -o ./udp-server/server ./udp-server/main.go

udp-client:
	go build -o ./udp-client/client ./udp-client/main.go

clean:
	rm -f ./tcp-server/server 
	rm -f ./tcp-client-single/client 
	rm -f ./tcp-client-multi/client 
	rm -f ./udp-server/server
	rm -f ./udp-client/client
