SHELL=PATH='$(PATH)' /bin/sh

GOBUILD=CGO_ENABLED=0 go build -ldflags '-w -s'

PLATFORM := $(shell uname -o)

COMMIT := $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := github.com/realbmail/web-client/utils
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME)

NAME := web_mta.exe
OS := windows

ifeq ($(PLATFORM), Msys)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else
	INCLUDE := $(GOPATH)
	NAME=web_mta
	OS=linux
endif

# enable second expansion
.SECONDEXPANSION:

.PHONY: all

BINDIR=$(INCLUDE)/bin

all:   build

build:
	GOOS=$(OS) GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)

target:=mac

mac:
	GOOS=darwin go build -ldflags '-w -s' -o $(NAME).mac  -ldflags="$(LD_FLAGS)" ./cmd/*.go
arm:
	CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARM=7 GOARCH=arm64 go build -ldflags '-w -s' -o $(BINDIR)/$(NAME).arm  -ldflags="$(LD_FLAGS)" ./cmd/*.go
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o $(NAME).lnx  -ldflags="$(LD_FLAGS)" ./cmd/*.go
win:
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o $(BINDIR)/$(NAME).exe  -ldflags="$(LD_FLAGS)" ./cmd/*.go

clean:
	rm $(BINDIR)/$(NAME)
