# Config
DEST   = out
PREFIX = /usr/local

# Commands
GO      = go
GOBUILD = $(GO) build
GOLIST  = $(GO) list
GOTEST  = $(GO) test
H2M     = help2man
SVU     = svu

# Target
LISTTPL = {{join .GoFiles " "}} {{join .EmbedFiles " "}}
FILES  != $(GOLIST) -f '$(LISTTPL)'
IMPORT != $(GOLIST) -f '{{.ImportPath}}'
NAME    = $(notdir $(IMPORT))
BIN     = $(DEST)/$(NAME)

# Docs
MAN   = $(DEST)/$(NAME).1
DESC != head -n1 usage.txt | sed 's/\.$$//'

# Build/test flags
VERSION   != $(SVU) current 2>/dev/null
TESTFLAGS  =
BUILDFLAGS = -trimpath
LDFLAGS    = -X 'main.version=$(VERSION)'

ifeq ($(VERSION),)
	VERSION = v0.0.0
endif
ifeq ($(VERSION),v0.0.0)
	VERSION := $(VERSION)-dev
endif

-include config.mk

GOTEST  := $(GOTEST) $(TESTFLAGS)
GOBUILD := $(GOBUILD) $(BUILDFLAGS) \
		   $(if $(LDFLAGS),-ldflags "$(LDFLAGS)")

all: $(BIN) $(MAN)

test:
	$(GOTEST)

install: install-bin install-man

install-bin: $(BIN)
	@mkdir -p $(PREFIX)/bin
	install -Dm755 -- $< $(PREFIX)/bin/$(<F)

install-man: $(MAN)
	@mkdir -p $(PREFIX)/share/man/man1
	install -Dm644 -- $< $(PREFIX)/share/man/man1/$(<F)

$(BIN): $(FILES)
	@mkdir -p $(@D)
	$(GOBUILD) -o $@

$(MAN): $(BIN)
	@mkdir -p $(@D)
	$(H2M) --output=$@ --name="$(DESC)" $<

show-%:
	@echo '$* = $($*)'

.PHONY: all test install install-bin install-man
