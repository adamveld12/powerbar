app := powerbar
PKGS := $(shell go list ./... | grep -v vendor)
OUTDIR := .bin
BINARY := $(OUTDIR)/$(app)
GOBIN := $(GOPATH)/bin
LINTBIN := $(GOBIN)/golangci-lint
clean_list = $(OUTDIR)

.PHONY: build clean dev test publish lint install

default: lint test build dev

$(OUTDIR):
	@mkdir .bin

$(BINARY): . $(OUTDIR)
	go build -o $@ $<

dev: clean $(BINARY)
	LB_DEBUG=true $(BINARY) $(FPATH)

build: $(BINARY)

clean: $(clean_list)
	rm -rf $<

install: build
	cp $(BINARY) $${HOME}/.config/waybar/


$(LINTBIN):
	@GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

lint: $(LINTBIN)
	go mod tidy -v
	$(LINTBIN) run -p bugs -p format -p performance -p unused

test: lint
	go test -v -cover $(PKGS)
