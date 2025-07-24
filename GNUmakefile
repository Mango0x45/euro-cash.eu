GO := go

cssfiles  := $(shell find static -name '*.css' -not -name '*.min.css')
cssfiles  := $(cssfiles:.css=.min.css)
gofiles   := $(shell find main.go src -name '*.go')
sqlfiles  := $(shell find src/dbx/sql -name '*.sql')
templates := $(shell find src/templates -name '*.tmpl')

exttmpl := $(wildcard cmd/exttmpl/*.go)

all: euro-cash.eu exttmpl

euro-cash.eu: $(cssfiles) $(templates) $(gofiles) $(sqlfiles)
	$(GO) build

po: exttmpl
	find . -name '*.html.tmpl' -exec ./exttmpl -out po/templates.pot {} +
	for bcp in en en-US nl;                                                     \
	do                                                                          \
		mkdir -p "po/$$bcp";                                                    \
		msgmerge --update "po/$$bcp/messages.po" po/templates.pot;              \
	done

exttmpl: $(exttmpl)
	$(GO) build ./cmd/exttmpl

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

.PHONY: clean po release