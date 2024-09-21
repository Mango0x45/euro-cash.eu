templates = $(shell find src/templates -name '*.tmpl')
gofiles = $(shell find main.go src -name '*.go')

exttmpl = $(wildcard cmd/exttmpl/*.go)
mfmt = $(wildcard cmd/mfmt/*.go)

all: euro-cash.eu exttmpl mfmt

euro-cash.eu: $(templates) $(gofiles)
	go build

# Generating translations is rather slow; so don’t do that by default
all-i18n: exttmpl
	go generate ./src
	find . -name out.gotext.json | mcp -b sed s/out/messages/
	go build

exttmpl: $(exttmpl)
	go build ./cmd/exttmpl

mfmt: $(mfmt)
	go build ./cmd/mfmt

watch:
	ls euro-cash.eu | entr -r ./euro-cash.eu -no-email -port $${PORT:-8080}

# Build a release tarball for easy deployment
# TODO: Minify CSS
release: all-i18n
	[ -n "$$GOOS" -a -n "$$GOARCH" ]
	tar -cf euro-cash.eu-$$GOOS-$$GOARCH.tar.gz euro-cash.eu data/ static/

clean:
	find . -type f \( \
		-name euro-cash.eu \
		-or -name exttmpl \
		-or -name mfmt \
		-or -name '*.min.css' \
		-or -name '*.tar.gz' \
	\) -delete
