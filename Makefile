lamport: lamport.go
	go build

clean:
	rm -rf lamport

all: lamport

.PHONY: all clean
