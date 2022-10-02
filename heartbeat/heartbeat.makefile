HEARTBEAT_BINARY=heartbeatApp
HEARTBEAT_SERVICE=heartbeat
PWD=$(pwd)

.PHONY: clean
clean:
	rm -rf ${HEARTBEAT_BINARY}

.PHONY: all
all: clean
	@echo "Building ${HEARTBEAT_BINARY} binary..." 
	cd ./${HEARTBEAT_SERVICE} && env GOOS=linux CGO_ENABLED=0 go build -o ${HEARTBEAT_BINARY} .
	@echo "Done!" 
