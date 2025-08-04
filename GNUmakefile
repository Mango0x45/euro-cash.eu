GO   := go
PORT := 8080

cssfiles  := $(shell find static -name '*.css' -not -name '*.min.css')
cssfiles  := $(cssfiles:.css=.min.css)
gofiles   := $(shell find main.go src -name '*.go')
sqlfiles  := $(shell find src/dbx/sql -name '*.sql')
templates := $(shell find src/templates -name '*.tmpl')

extpo   := $(wildcard cmd/extpo/*.go)
extwiki := $(wildcard cmd/extwiki/*.go)

ENABLED_LANGUAGES := $(shell ./aux/enabled-languages)

all: euro-cash.eu extpo

euro-cash.eu: $(cssfiles) $(templates) $(gofiles) $(sqlfiles)
	$(GO) build

extract: extpo extwiki
	find . -name '*.go' -exec xgotext --foreign-user -o po/backend.pot {} +
	find . -name '*.html.tmpl' -exec ./extwiki {} +                             \
		| gofmt >src/wikipedia/links.gen.go
	find . -name '*.html.tmpl' -exec ./extpo {} +                               \
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

extpo: $(extpo)
	$(GO) build ./cmd/extpo

extwiki: $(extwiki)
	$(GO) build ./cmd/extwiki

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
		-or -name extpo                                                         \
		-or -name extwiki                                                       \
		-or -name '*.min.css'                                                   \
		-or -name '*.tar.gz'                                                    \
	\) -delete

debug:
	./euro-cash.eu -debug -no-email -db-name :memory: -port $(PORT)

.PHONY: clean debug extract po release
