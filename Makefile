## ----- Variables -----
PKG_NAME = $(shell basename "$$(pwd)")
ifeq ($(shell ls -1 go.mod 2> /dev/null),go.mod) # use module name from go.mod, if applicable
	PKG_NAME = $(shell basename "$$(cat go.mod | grep module | awk '{print $$2}')")
endif

VERSION = $(shell git describe --tags | cut -c 2-)

## Directory of the 'main' package.
MAINDIR = "."
## Output directory to place artifacts from 'build' and 'build-all'.
OUTDIR  = "."

## Enable Go modules for this project.
MODULES = true
## Enable goreleaser for this project.
GORELEASER = false
## Enable git-secret for this project.
SECRETS = false

## Custom Go linker flags: (DISABLED)
# LDFLAGS = -X github.com/stevenxie/$(PKG_NAME)/cmd.Version=$(VERSION)


## Source configs:
SRC_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SRC_PKGS = $(shell go list ./... | grep -v /vendor/)

## Testing configs:
TEST_TIMEOUT = 20s
COVER_OUT = coverage.out



## ------ Commands (targets) -----
.PHONY: default setup init

## Default target when no arguments are given to make (build and run program).
default: build-run

## Sets up this project on a new device.
setup: hooks-setup
	@if [ "$(SECRETS)" == true ]; then $(SECRETS_REVEAL_CM); fi
	@if [ "$(MODULES)" == true ]; \
	 then $(DL_CMD); \
	 else $(GET_CMD); \
	 fi

## Initializes this project from scratch.
## Variables: MODPATH
init: mod-init secrets-init goreleaser-init


## [Git, git-secret]
.PHONY: hooks-setup secrets-hide secrets-reveal

## Configure Git to use .githooks (for shared githooks).
hooks-setup:
	@echo "Configuring githooks..."
	@git config core.hooksPath .githooks && echo "done"

## Initialize git-secret.
secrets-init:
	@if [ "$(SECRETS)" == true ]; then \
	   echo "Initializing git-secret..." && \
	   git-secret init; \
	 fi

## Hide modified secret files using git-secret.
secrets-hide:
	@echo "Hiding modified secret files..."
	@git secret hide -m

## Reveal files hidden by git-secret.
SECRETS_REVEAL_CM = git secret reveal
secrets-reveal:
	@echo "Revealing hidden secret files..."
	@$(SECRETS_REVEAL_CM)


## [Go: modules]
.PHONY: init verify dl vendor tidy update fix

## Initializes a Go module in the current directory.
## Variables: MODPATH (module source path)
MODPATH =
mod-init:
	@if [ "$(MODULE)" == true ]; then \
	   echo "Initializing Go module..." && \
	   go mod init $(MODPATH); \
	 fi

## Verifies that Go module dependencies are satisfied.
VERIFY_CMD = echo "Verifying Go module dependencies..." && go mod verify
verify:
	@$(VERIFY_CMD)

## Downloads Go module dependencies.
DL_CMD = echo "Downloading Go module dependencies..." && \
           go mod download && echo "done"
dl:
	@$(DL_CMD)

## Vendors Go module dependencies.
vendor:
	@echo "Vendoring Go module dependencies..."
	@go mod vendor && echo "done"

## Tidies Go module dependencies.
tidy:
	@echo "Tidying Go module dependencies..."
	@go mod tidy && echo "done"

## Installs and updates package dependencies.
## Variables:
##   UMODE (Update Mode, choose between 'patch' and 'minor').
UMODE =
update:
	@echo 'Updating module dependencies with "go get -u"...'
	@go get -u $(UMODE) && echo "done"

## Fixes deprecated Go code using "go fix", by rewriting old APIS to use
## newer ones.
fix:
	@echo 'Fixing deprecated Go code with "go fix"... '
	@go fix && echo "done"


## [Go: legacy setup]
.PHONY: get

## Downloads and installs all subpackages (legacy).
GET_CMD = echo "Installing dependencies... " && \
          go get ./... && echo "done"
get:
	@$(GET_CMD)


## [Go: setup, running]
.PHONY: build build-all build-run run clean install

## Runs the built program.
## Sources .env.sh if it exists.
## Variables: SRCENV (boolean which determines whether or not to check and
##            source .env.sh)
SRCENV = true
OUTPATH = $(OUTDIR)/$(PKG_NAME)
RUN_CMD = \
	if [ -f ".env.sh" ] && [ "$(SRCENV)" == true ]; then \
	  echo 'Configuring environment variables by sourcing ".env.sh"...' && \
	  . .env.sh && \
	  printf "done\n\n"; \
	fi; \
	if [ -f "$(OUTPATH)" ]; then \
	  echo 'Running "$(PKG_NAME)"...' && \
	  ./$(OUTPATH); \
	else \
	  echo 'run: could not find program "$(OUTPATH)".' >&2; \
	  exit 1; \
	fi
run:
	@$(RUN_CMD)

## Builds (compiles) the program for this system.
## Variables:
##   - OUTDIR (output directory to place built binaries)
##   - MAINDIR (directory of the main package)
##   - BUILDARGS (additional arguments to pass to "go build")
BUILDARGS =
BUILD_CMD = \
	echo 'Building "$(PKG_NAME)" for this system...' && \
	go build \
	  -o "$$(echo $(OUTDIR) | tr -s '/')/$(PKG_NAME)" \
	  -ldflags "$(LDFLAGS)" \
	  $(BUILDARGS) $(MAINDIR) && \
	echo "done"
build:
	@$(BUILD_CMD)

## Builds (cross-compiles) the program for all systems.
## Variables:
##   - OUTDIR (output path to place built binaries)
##   - MAINDIR (directory of the main package)
##   - BUILDARGS (additional arguments to pass to "go build")
build-all:
	@echo 'Building "$(PKG_NAME)" for all systems:'
	@for GOOS in darwin linux windows; do \
	   for GOARCH in amd64 386; do \
	     printf "Building GOOS=$$GOOS GOARCH=$$GOARCH... " && \
	     OUTNAME="$(PKG_NAME)-$$GOOS-$$GOARCH"; \
	     if [ $$GOOS == windows ]; then \
	       OUTNAME="$$OUTNAME.exe"; \
	     fi; \
	     GOBUILD_OUT="$$(GOOS=$$GOOS GOARCH=$$GOARCH && \
	       go build \
	         -o "$$(echo $(OUTDIR) | tr -s '/')/$$OUTNAME" \
	         -ldflags "$(LDFLAGS)" \
	         $(BUILDARGS) $(MAINDIR) 2>&1)"; \
	     if [ -n "$$GOBUILD_OUT" ]; then \
	       printf "\nError during build:\n" >&2 && \
	        echo "$$GOBUILD_OUT" >&2 && \
	        exit 1; \
	     else printf "\tdone\n"; \
	     fi; \
	   done; \
	 done

## Builds (compiles) the program for this system, and runs it.
## Sources .env.sh before running, if it exists.
build-run:
	@$(BUILD_CMD) && echo "" && $(RUN_CMD)

## Cleans build artifacts (executables, object files, etc.).
clean:
	@echo 'Cleaning build artifacts with "go clean"...'
	@go clean && echo "done"

## Installs the program using "go install".
install:
	@echo 'Installing program using "go install"... '
	@go install && echo "done"


## [Go: code checking]
.PHONY: fmt lint vet check

## Formats the source code using "gofmt".
FMT_CMD = \
	if ! command -v gofmt > /dev/null; then \
	  echo '"gofmt" is required to format source code.'; \
	else \
	  echo 'Formatting source code using "gofmt"...' && \
	  gofmt -l -s -w . && echo "done"; \
	fi
fmt:
	@$(FMT_CMD)

## Lints the source code using "golint".
LINT_CMD = \
	if ! command -v golint > /dev/null; then \
	  echo '"golint" is required to lint soure code.' >&2; \
	else \
	  echo 'Formatting source code using "golint"...' && \
	  golint ./... && echo "done"; \
	fi
lint:
	@$(LINT_CMD)

## Checks for suspicious code using "go vet".
VET_CMD = echo 'Checking for suspicious code using "go vet"...' && \
	      go vet && echo "done"
vet:
	@$(VET_CMD)

## Checks for formatting, linting, and suspicious code.
CHECK_CMD = $(FMT_CMD) && echo "" && $(LINT_CMD) && echo "" && $(VET_CMD)
check:
	@$(CHECK_CMD)


## [Go: testing]
.PHONY: test test-v test-race test-race-v bench bench-v

TEST_CMD = go test ./... -coverprofile=$(COVER_OUT) \
                         -covermode=atomic \
                         -timeout=$(TEST_TIMEOUT)
test:
	@echo "Testing:"
	@$(TEST_CMD)
test-v:
	@echo "Testing (verbose):"
	@$(TEST_CMD) -v

TEST_CMD_RACE = $(TEST_CMD) -race
test-race:
	@echo "Testing (race):"
	@$(TEST_CMD_RACE)
test-race-v:
	@echo "Testing (race, verbose):"
	@$(TEST_CMD_RACE) -v

BENCH_CMD = $(TEST_CMD) ./... -run=^$$ -bench=. -benchmem
bench:
	@echo "Benchmarking:"
	@$(BENCH_CMD)
bench-v:
	@echo "Benchmarking (verbose):"
	@$(BENCH_CMD) -v


## [Go: reviewing]
.PHONY: review review-race review-bench
__review_base:
	@$(VERIFY_CMD) && echo "" && $(CHECK_CMD) && echo ""

## Formats, checks, and tests the code.
review: __review_base test
review-v: __review_base test-v

## Like "review", but tests for race conditions.
review-race: __review_base test-race
review-race-v: __review_base test-race-v

## Like "review-race", but includes benchmarks.
review-bench: review-race bench
review-bench-v: review-race bench-v


## [Goreleaser]
.PHONY: goreleaser-init release

goreleaser-init:
	@if [ "$(GORELEASER)" == true ]; then \
	   echo "Initializing goreleaser..." && \
	   goreleaser init; \
	 fi

release:
	@echo "Releasing with 'goreleaser'..." && goreleaser --rm-dist

snapshot:
	@echo "Making snapshot with 'goreleaser'..." && \
	   goreleaser --snapshot --rm-dist
