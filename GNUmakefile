GO := go

cssfiles  := $(shell find static -name '*.css' -not -name '*.min.css')
cssfiles  := $(cssfiles:.css=.min.css)
gofiles   := $(shell find main.go src -name '*.go')
sqlfiles  := $(shell find src/dbx/sql -name '*.sql')
templates := $(shell find src/templates -name '*.tmpl')

exttmpl := $(wildcard cmd/exttmpl/*.go)

ENABLED_LANGUAGES := $(shell ./enabled-languages)

all: euro-cash.eu exttmpl

euro-cash.eu: $(cssfiles) $(templates) $(gofiles) $(sqlfiles)
	$(GO) build

extract: exttmpl
	find . -name '*.go' -exec xgotext --foreign-user -o po/backend.pot {} +
	find . -name '*.html.tmpl' -exec ./exttmpl {} + \
		| msgcat po/backend.pot - -o po/messages.pot
	for bcp in $(ENABLED_LANGUAGES);                                            \
	do                                                                          \
		dir="po/$$bcp";                                                         \
		if [ ! -d "$$dir" ];                                                    \
		then                                                                    \
			mkdir -p "$$dir";                                                   \
			msginit -i po/messages.pot -o "$$dir/messages.po" -l$$bcp.UTF-8     \
				--no-translator;                                                \
		fi;                                                                     \
		msgmerge -UN "po/$$bcp/messages.po" po/messages.pot;                    \
	done
	find po -name '*~' -delete

po:
	for po in po/*/*.po;                                                        \
	do                                                                          \
		msgfmt "$$po" -o "$${po%.*}.mo";                                        \
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

.PHONY: clean extract po release