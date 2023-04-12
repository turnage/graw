DOCKER_COMPOSE_TEST := docker-compose -f dev/test.yml
TEST_SERVICE_NAME := test_graw

ifdef TEST_RUN
	TESTRUN := -run ${TEST_RUN}
endif

GOPACKAGES := $(shell go list ./... | egrep -v '/vendor|github.com/mix/graw$$')
TEST_MODULES ?= $(GOPACKAGES)


test: # run unit tests
	${DOCKER_COMPOSE_TEST} rm --force || true
	${DOCKER_COMPOSE_TEST} run ${TEST_SERVICE_NAME}
	${DOCKER_COMPOSE_TEST} down

test-direct: # [INTERNAL]
	go test -p 1 -v -race  -coverprofile=$(COVERAGE_FILE) $(TEST_MODULES) $(TESTRUN)

lint: # Run go lint
	${DOCKER_COMPOSE_TEST} run test_graw bash -c "GOGC=50 make -e lint-direct"

lint-direct: # [INTERNAL]
	@golangci-lint run

stop: # stop services
	${DOCKER_COMPOSE_TEST} down || true

update: stop
	${DOCKER_COMPOSE_TEST} rm --force ${TEST_SERVICE_NAME}
	${DOCKER_COMPOSE_TEST} build ${TEST_SERVICE_NAME}
