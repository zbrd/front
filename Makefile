.DEFAULT_GOAL := all

# Config
DEST   = out
PREFIX = /usr/local
VERVAR = main.version
SRCS   = Go Embed

# Commands
GIT     = git
GO      = go
GOBUILD = $(GO) build
GOLIST  = $(GO) list
GOLF    = $(GOLIST) -f
GOTEST  = $(GO) test
INSTALL = install
H2M     = help2man

# Main packages
PKGS != $(GOLF) '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./...
BINS  = $(addprefix $(DEST)/,$(notdir $(PKGS)))
MANS  = $(addsuffix .1,$(BINS))

# Build/test flags
VERSION    != $(GIT) tag --sort=-v:refname 2>/dev/null | head -n1
LDFLAGS     = $(if $(VERSION),-X '$(VERVAR)=$(VERSION)')
BUILDFLAGS  = -trimpath $(if $(LDFLAGS),-ldflags "$(LDFLAGS)")
TESTFLAGS   = -v

# Functions
getrange = {{range .$(1)Files}}{{$$.Dir}}/{{.}} {{end}}
filesfmt = $(foreach t,$(SRCS),$(call getrange,$(t)))
getpreqs = $(shell $(GOLF) '$(filesfmt)' $(1))
getdir   = $(shell $(GOLF) '{{.Dir}}' $(1))
getusage = $(shell test -f $(1) && head -n1 $(1) \
		   | sed -e 's/^[A-Z]/\L&/' -e 's/\.$$//')

# Rule generator
define mkpkg =
$(1)_USAGE = $(call getusage,$(call getdir,$(2))/usage.txt)

$$(DEST)/$(1): $(call getpreqs,$(2))
	@mkdir -p $$(@D)
	$$(GOBUILD) $$(BUILDFLAGS) -o $$@ $(2)

$$(DEST)/$(1).1: USAGE ?= $$($(1)_USAGE)
$$(DEST)/$(1).1: $(call getpreqs,$(2))
	@mkdir -p $$(@D)
	$$(H2M) --output=$$@ --name='$$(USAGE)' $$(DEST)/$(1)
endef

-include config.mk

all: $(BINS) $(MANS)

test:
	$(GOTEST) $(TESTFLAGS) ./...

install: $(BINS) $(MANS)
	$(INSTALL) -Dm755 $(BINS) $(PREFIX)/bin/
	$(INSTALL) -Dm644 $(MANS) $(PREFIX)/share/man/man1/

clean:
	-rm $(BINS)

$(foreach p,$(PKGS),$(eval $(call mkpkg,$(notdir $(p)),$(p))))

show-%:
	@echo '$* = $($*)'

.PHONY: all test install clean
