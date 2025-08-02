package wikipedia

type APIResponse struct {
	Continue *struct {
		Continue   string `json:"continue"`
		LlContinue string `json:"llcontinue"`
	} `json:"continue"`
	Query struct {
		Pages []struct {
			Title     string `json:"title"`
			LangLinks *[]struct {
				Lang  string `json:"lang"`
				Title string `json:"title"`
			} `json:"langlinks"`
		} `json:"pages"`
	} `json:"query"`
}
