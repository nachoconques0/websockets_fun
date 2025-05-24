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