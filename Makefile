HEARTBEAT_BINARY=heartbeatApp
MAIL_BINARY=mailerApp

## up: starts all containers in the background without forcing build
.PHONY: up
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
.PHONY: up_build
up_build: build_heartbeat 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
.PHONY: down 
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!" 

## build_heartbeat: builds the heartbeatApp binary as a linux executable
.PHONY: build_heartbeat 
build_heartbeat:
	@echo "Building ${HEARTBEAT_BINARY} binary..." 
	cd ./heartbeat && env GOOS=linux CGO_ENABLED=0 go build -o ${HEARTBEAT_BINARY} .
	@echo "Done!" 

## build_mail: builds the mail binary as a linux executable 
.PHONY: build_mail
build_mail:
	@echo "Building mail binary..."
	# cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} .
	@echo "Done!"
