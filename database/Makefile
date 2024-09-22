SERVER_APP_NAME=spider-server
CLI_APP_NAME=spider-cli

build-server:
	go build -o ${SERVER_APP_NAME} cmd/server/main.go

build-cli:
	go build -o ${CLI_APP_NAME} cmd/cli/main.go

run-server: build-server
	./${SERVER_APP_NAME}

run-server-with-config: build-server
	CONFIG_FILE_NAME=config.yml ./${SERVER_APP_NAME}

run-cli: build-cli
	./${CLI_APP_NAME} $(ARGS)

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out