# Generating translations is rather slow; so don’t do that by default
all:
	TEMPL_EXPERIMENT=rawgo go generate ./templates
	go build

all-i18n:
	TEMPL_EXPERIMENT=rawgo go generate ./templates ./i18n
	find . -name out.gotext.json | mcp -b sed s/out/messages/
	go build

watch:
	ls euro-cash.eu | entr -r ./euro-cash.eu

.PHONY: watch
