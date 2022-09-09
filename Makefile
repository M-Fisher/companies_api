GOCI_VERSION=v1.49.0
NAMESPACE = companies-service
VERSION ?= $(shell git describe --tags || git rev-parse --short HEAD)

_COMPOSE=docker-compose -f dev/docker-compose.yaml --project-name companies-service
_COMPOSE_BUILD=${_COMPOSE} --profile app build

## CI Targets
# Docker images building
@build: build-app build-migrations
	$(info [MAKE] Bulding production image)

@test: docker-build-ci-image
	$(info [MAKE] Running tests)
	@docker run --rm \
		${NAMESPACE}/ci:latest \
		make test

@lint: docker-build-ci-image
	$(info [MAKE] Running linters)
	@docker run --rm \
		${NAMESPACE}/ci:latest \
		make lint

# App binary building
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -v -a -o output main.go

build-app:
	@docker build \
		--tag ${NAMESPACE}/app:latest \
		--tag ${NAMESPACE}/app:${VERSION} \
		--file dev/Dockerfile \
		.

build-migrations:
	@docker build \
    		--file dev/migrations/Dockerfile \
    		--tag ${NAMESPACE}/migrations:${VERSION} \
			--tag ${NAMESPACE}/migrations:latest \
    		.

docker-build-ci-image:
	@docker build \
		--tag ${NAMESPACE}/ci:latest \
		--tag ${NAMESPACE}/ci:${VERSION} \
		--file dev/Dockerfile \
		--target ci \
		.

setup: GOLANGCI_INSTALLER=https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
setup:
	curl -sSfL ${GOLANGCI_INSTALLER} | sh -s -- -b ${GOPATH}/bin ${GOCI_VERSION}

## Local targets
dev-up-app:
	${_COMPOSE_BUILD}
	${_COMPOSE} --profile app up

dev-down-app:
	${_COMPOSE} down

dev-down:
	${_COMPOSE} down --remove-orphans

dev-clean:
	${_COMPOSE} down -v --rmi all

dev-restart: dev-down dev-up-app

dev-local-env:
	${_COMPOSE} --profile dependencies up

dev-migrate:
	${_COMPOSE} run --rm --service-ports $(NAMESPACE)-migrate

test:
	CGO_ENABLED=1 go test -cover -race -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

local-run:
	SERVICE_CONF=${PWD}/dev/local.env go run main.go server

generate-mocks: 
	mockery --recursive --name=EventsService --inpackage
	mockery --recursive --name=CompaniesService --inpackage
	mockery --recursive --name=IPDataProvider --inpackage
	mockery --recursive --name=AuthService --inpackage
	mockery --recursive --name=CompaniesQueries --inpackage
	
