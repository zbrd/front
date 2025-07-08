# Config
DEST   = bin
PREFIX = /usr/local

# Commands
GO      = go
GOBUILD = $(GO) build
GOLIST  = $(GO) list
GOTEST  = $(GO) test

# Target
LISTTPL = {{join .GoFiles " "}} {{join .EmbedFiles " "}}
FILES  != $(GOLIST) -f '$(LISTTPL)'
IMPORT != $(GOLIST) -f '{{.ImportPath}}'
NAME    = $(notdir $(IMPORT))
BIN     = $(DEST)/$(NAME)

# Build/test flags
TESTFLAGS  =
BUILDFLAGS = -trimpath
LDFLAGS    =

-include config.mk

GOTEST  := $(GOTEST) $(TESTFLAGS)
GOBUILD := $(GOBUILD) $(BUILDFLAGS) \
		   $(if $(LDFLAGS),-ldflags "$(LDFLAGS)")

all: $(BIN)

test:
	$(GOTEST)

install: $(BIN)
	@mkdir -p $(PREFIX)/bin
	install -Dm755 -- $< $(PREFIX)/bin/$(<F)

$(BIN): $(FILES)
	@mkdir -p $(@D)
	$(GOBUILD) -o $@

show-%:
	@echo '$* = $($*)'

.PHONY: all test install
