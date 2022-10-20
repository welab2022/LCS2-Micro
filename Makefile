AUTHENTICATION_BINARY=authenticationApp
MAIL_BINARY=mailerApp

TEST_DIR=test
TEST_REPORT=test/report/report.html
## up: starts all containers in the background without forcing build
.PHONY: up
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
.PHONY: up_build
up_build: build_auth
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
	@echo

## build_heartbeat: builds the heartbeatApp binary as a linux executable
.PHONY: build_auth
build_auth: clean_auth
	@echo
	@echo "Building ${AUTHENTICATION_BINARY} binary..." 
	cd ./authentication && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATION_BINARY} ./cmd/api
	@echo "Done!"
	@echo

## build_mail: builds the mail binary as a linux executable 
.PHONY: build_mail
build_mail:
	@echo "Building mail binary..."
	# cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} .
	@echo "Done!"


#############
## test_all: tests all the services
.PHONY: test_all
test_all:
	@echo
	@echo "Testing all the services..."
	pytest -v ${TEST_DIR} --html=${TEST_REPORT}
	open ${TEST_REPORT}
	@echo "Done!"
	@echo

#############
## clean_auth: delete all objects and binaries of Authentication service
.PHONY: clean_auth
clean_auth:
	@echo
	@echo "Cleaning Authentication service binaries..."
	cd ./authentication && rm -rf ${AUTHENTICATION_BINARY}
	@echo "Done!"
	@echo 

## clean_all: delete all objects and binaries of all the services
.PHONY: clean
clean: clean_auth
	@echo
	@echo "Cleaning up..."
	@echo "Done!"
	@echo 
