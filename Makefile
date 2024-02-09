VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=$(VERSION)
else
	TAG=$(version)
endif
BIN_DIR=bin

ifneq ($(verbose),)
	TEST_ARGS += -v
endif

ifeq ($(short),true)
	TEST_ARGS += -short
endif

PACKAGE=github.com/gre-ory/games-go

LDFLAGS = -X 'main.version=$(TAG)'

# run 'make Q="" <rule>' to enable verbosity
Q := @

.PHONY:	all build test run

build:
	@echo " ----- build -----"
	$(Q)CGO_ENABLED=1 GOOS=linux go build $(GO_BUILD_FLAGS) -ldflags "${LDFLAGS}" -o ${BIN_DIR}/server ${PACKAGE}/server
run: build
	@echo " ----- run -----"
	@./scripts/run
test:
	$(Q) go test -race ./...
