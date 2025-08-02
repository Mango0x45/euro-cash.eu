package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	defaultLocale string
	titlemap      = make(map[string]map[string]string)
)

func Init(locale string) {
	defaultLocale = locale
	base := fmt.Sprintf("https://%s.wikipedia.org/w/api.php", defaultLocale)
	u, err := url.Parse(base)
	if err != nil {
		log.Println(err)
		return
	}

	var resp APIResponse
	titles := strings.Join(extractedTitles[:], "|")

	q := u.Query()
	q.Set("action", "query")
	q.Set("format", "json")
	q.Set("prop", "langlinks")
	q.Set("titles", titles)
	q.Set("formatversion", "2")
	q.Set("lllimit", "max")

	for {
		if resp.Continue != nil {
			q.Set("continue", resp.Continue.Continue)
			q.Set("llcontinue", resp.Continue.LlContinue)
		}
		u.RawQuery = q.Encode()

		respjson, err := http.Get(u.String())
		if err != nil {
			log.Println(err)
			return
		}
		defer respjson.Body.Close()

		body, err := io.ReadAll(respjson.Body)
		if err != nil {
			log.Println(err)
			return
		}

		resp = APIResponse{}
		if err = json.Unmarshal(body, &resp); err != nil {
			log.Println(err)
			return
		}

		for _, page := range resp.Query.Pages {
			if page.LangLinks == nil {
				continue
			}

			t := url.PathEscape(page.Title)
			if _, ok := titlemap[t]; !ok {
				titlemap[t] = make(map[string]string)
			}

			for _, ll := range *page.LangLinks {
				titlemap[t][ll.Lang] = url.PathEscape(ll.Title)
			}
		}

		if resp.Continue == nil {
			break
		}
	}
}

func Url(title, locale string) string {
	title = url.PathEscape(title)
	t, ok := titlemap[title][locale]
	if !ok {
		return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s",
			defaultLocale, title)
	}
	return fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", locale, t)
}
