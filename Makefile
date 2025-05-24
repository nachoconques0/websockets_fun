REDIS_HOST=localhost
REDIS_STREAM=local-q

.PHONY: mod
## Install project dependencies using go mod. Usage 'make mod'
mod:
	@go mod tidy
	@go mod vendor

.PHONY: mock
## Generate mock files. Usage: 'make mock'
mock: ; $(info Generating mock files)
	chmod +x generate-mocks.sh
	@./generate-mocks.sh

.PHONY: run
## Run service. Usage: 'make run'
run: ; $(info Starting svc...)
	go run --tags dev ./main.go

.PHONY: test
## Run tests. Usage: 'make test' Options: path=./some-path/... [and/or] func=TestFunctionName
test: ; $(info running tests…) @
	@docker compose -f ./docker-compose.yml up --quiet-pull --force-recreate -d --wait;
	@if [ -z $(path) ]; then \
		path='./...'; \
	else \
		path=$(path); \
	fi; \
	if [ ! -d "coverage" ]; then \
		mkdir coverage; \
	fi; \
	if [ -z $(func) ]; then \
		go test -v -failfast -covermode=count -coverprofile=./coverage/coverage.out $$path; \
	else \
		go test -v -failfast -covermode=count -coverprofile=./coverage/coverage.out -run $$func $$path; \
	fi;
	docker compose -f ./docker-compose.yml down;
