cssfiles  := $(shell find static -name '*.css' -not -name '*.min.css')
cssfiles  := $(cssfiles:.css=.min.css)
gofiles   := $(shell find main.go src -name '*.go')
templates := $(shell find src/templates -name '*.tmpl')

exttmpl := $(wildcard cmd/exttmpl/*.go)
mfmt    := $(wildcard cmd/mfmt/*.go)

all: euro-cash.eu exttmpl mfmt

euro-cash.eu: $(cssfiles) $(templates) $(gofiles)
	go build

all-i18n: exttmpl
	go generate ./src
	find . -name out.gotext.json | mcp -b sed s/out/messages/
	go build

exttmpl: $(exttmpl)
	go build ./cmd/exttmpl

mfmt: $(mfmt)
	go build ./cmd/mfmt

%.min.css: %.css
	lightningcss -m $< -o $@

watch:
	ls euro-cash.eu | entr -r ./euro-cash.eu -no-email -port $${PORT:-8080}

release: all-i18n
	[ -n "$$GOOS" -a -n "$$GOARCH" ]
	find data static -type f \(                                                 \
		-not -name '*.css'                                                      \
		-or -name '*.min.css'                                                   \
	\) -exec tar -cf euro-cash.eu-$$GOOS-$$GOARCH.tar.gz euro-cash.eu {} +

clean:
	find . -type f \(                                                           \
		-name euro-cash.eu                                                      \
		-or -name exttmpl                                                       \
		-or -name mfmt                                                          \
		-or -name '*.min.css'                                                   \
		-or -name '*.tar.gz'                                                    \
	\) -delete
