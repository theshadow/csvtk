EXE = "csvtk"
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

BIN = $(GOPATH)/bin

GO = /usr/bin/env go
DEP = /usr/bin/env dep
FMT = $(GO) fmt
VET = $(GO) vet

$(EXE): vendor static-analysis ; $(info $(M) building $@…)
	$Q $(GO) build \
		-tags release \
		-ldflags '-X github.com/theshadow/csvtk/cmd.Version=$(VERSION) -X github.com/theshadow/csvtk/cmd.BuildDate=$(DATE)' \
		-o bin/csvtk main.go

.PHONY: build
build: $(EXE) ; $(info $(M) starting build…)

# Tools
GOLINT = $(BIN)/golint
$(GOLINT): ; $(info $(M) building golint…)
	$Q go get github.com/golang/lint/golint

# Depedencies
Gopkg.toml: ; $(info $(M) initializing dependencies…)
	$Q $(DEP) init
	@touch $@
Gopkg.lock: Gopkg.toml ; $(info $(M) updating dependencies…)
	$Q $(DEP) ensure -update
	@touch $@
vendor: Gopkg.lock ; $(info $(M) downloading dependencies…)
	$Q $(DEP) ensure

# Testing
.PHONY: test $(info $(M) running tests…)
test: vendor ;
	$(GO) test ./...

# Static Analysis
# See https://github.com/mre/awesome-static-analysis#Go for more options
.PHONY: static-analysis
static-analysis: vendor fmt vet lint

.PHONY: fmt
fmt: ; $(info $(M) running go fmt…)
	$Q $(GO) fmt ./...

# For now we'll just scan everything. Hopefully they follow the standard and exclude vendor/ when using '...'.
# (https://github.com/golang/lint/issues/320)
.PHONY: lint
lint: $(GOLINT) vendor ; $(info $(M) running golint…)
	$Q ret=0 && test -z "$$($(GOLINT) ./... | grep -v "vendor/" | tee /dev/stderr)" || ret=1 ; exit $$ret

.PHONY: vet
vet: vendor ; $(info $(M) running go vet…)
	$Q $(VET) ./...

# Misc
.PHONY: version
version: ; $(info $(M) running go fmt…) @ ## Run gofmt on all source files
	@echo $(VERSION)