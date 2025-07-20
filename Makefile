PREFIX = /usr/local

GO      = go
GLF     = go list -f
GIT     = git
INSTALL = install

GOFILES != $(GLF) '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...
EMFILES != $(GLF) '{{range .EmbedFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...
FILES    = $(GOFILES) $(EMFILES)

VERSION   != $(GIT) tag --sort=v:refname 2>/dev/null | head -n1
BUILDARGS  = -trimpath -ldflags "-X main.version=$(VERSION)"

ifeq ($(VERSION),)
	VERSION = v0.0.0-dev
endif

-include config.mk

all: out/front

test:
	$(GO) test ./...

install: out/front
	$(INSTALL) -Dm755 $< $(PREFIX)/bin/$(<F)

out/%: go.mod $(FILES)
	$(GO) build -o $@ $(BUILDARGS) ./cmd/$*

show-%:
	@echo $(*) = '$($*)'
