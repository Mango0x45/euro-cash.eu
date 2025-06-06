cssfiles  := $(shell find static -name '*.css' -not -name '*.min.css')
cssfiles  := $(cssfiles:.css=.min.css)
gofiles   := $(shell find main.go src -name '*.go')
templates := $(shell find src/templates -name '*.tmpl')

exttmpl := $(wildcard cmd/exttmpl/*.go)

all: euro-cash.eu exttmpl

euro-cash.eu: $(cssfiles) $(templates) $(gofiles)
	go build

all-i18n: exttmpl
	go generate ./src
	find . -name out.gotext.json | mcp -b sed s/out/messages/
	go build

exttmpl: $(exttmpl)
	go build ./cmd/exttmpl

%.min.css: %.css
	if command -v lightningcss >/dev/null;                                      \
	then                                                                        \
		lightningcss -m $< -o $@;                                               \
	else                                                                        \
		cp $< $@;                                                               \
	fi

clean:
	find . -type f \(                                                           \
		-name euro-cash.eu                                                      \
		-or -name exttmpl                                                       \
		-or -name '*.min.css'                                                   \
		-or -name '*.tar.gz'                                                    \
	\) -delete

.PHONY: all-i18n clean release