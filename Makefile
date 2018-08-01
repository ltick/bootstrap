default: all

init:
	@echo GOROOT=$(GOROOT)
	@echo GOPATH=$(GOPATH)
	@echo LIBRARY_PATH=$(LIBRARY_PATH)
	@echo PKG_CONFIG_PATH=$(PKG_CONFIG_PATH)
	@echo CGO_CFLAGS=$(CGO_CFLAGS)
	@echo CGO_LDFLAGS=$(CGO_LDFLAGS)

dummy-server: init
ifeq ($(DEBUG), TRUE)
	go build -race -ldflags=-s -o bin/dummy-server -v github.com/ltick/dummy/app/dummy-server
else
	go build -ldflags=-s -o bin/dummy-server -v github.com/ltick/dummy/app/dummy-server
endif
all: dummy-server

test: export PREFIX_PATH=$(CURDIR)
test: init
	go test -ldflags=-s -coverprofile cover.out -cover nebula/dummy/test/client/go

install: init
	go install -ldflags=-s -o bin/dummy-server -v nebula/dummy/app/dummy-server

clean: init
	go clean

.DEFAULT: install
