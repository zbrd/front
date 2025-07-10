# Config
DEST   = out
PREFIX = /usr/local

# Commands
GO      = go
GOBUILD = $(GO) build
GOLIST  = $(GO) list
GOTEST  = $(GO) test
H2M     = help2man
INSTALL = install
SVU     = svu

# Target
LISTTPL  = {{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}
LISTTPL += {{range .EmbedFiles}}{{$$.Dir}}/{{.}} {{end}}
FILES   != $(GOLIST) -f '$(LISTTPL)' ./...
IMPORT  != $(GOLIST) -f '{{.ImportPath}}' ./...
NAME     = $(notdir $(IMPORT))
BIN      = $(addprefix $(DEST)/,$(NAME))

# Docs
MAN   = $(addsuffix .1,$(BIN))
DESC != head -n1 usage.txt | sed 's/\.$$//'

# Build/test flags
VERSION    != $(SVU) current 2>/dev/null
TESTFLAGS   = -v
LDFLAGS     = -X 'main.version=$(VERSION)'
BUILDFLAGS  = -trimpath $(if $(LDFLAGS),-ldflags "$(LDFLAGS)")

ifeq ($(VERSION),)
	VERSION = v0.0.0
endif

ifeq ($(VERSION),v0.0.0)
	VERSION := $(VERSION)-dev
endif

-include config.mk

all: $(BIN) $(MAN)

test:
	$(GOTEST) $(TESTFLAGS)

install: install-bin install-man

install-bin: $(BIN)
	@mkdir -p $(PREFIX)/bin
	$(INSTALL) -m755 -Dt $(PREFIX)/bin -- $^

install-man: $(MAN)
	@mkdir -p $(PREFIX)/share/man/man1
	$(INSTALL) -m644 -Dt $(PREFIX)/share/man/man1 -- $^

$(BIN) &: $(FILES)
	@mkdir -p $(@D)
	$(GOBUILD) $(BUILDFLAGS) -o $(@D) ./...

$(MAN): $(BIN)
	@mkdir -p $(@D)
	$(H2M) --output=$@ --name="$(DESC)" $<

show-%:
	@echo '$* = $($*)'

.PHONY: all test install install-bin install-man
