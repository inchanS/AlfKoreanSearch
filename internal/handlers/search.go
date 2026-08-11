package handlers

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/inchanS/AlfKoreanSearch/internal/alfred"
	"github.com/inchanS/AlfKoreanSearch/internal/cache"
	"github.com/inchanS/AlfKoreanSearch/internal/httpx"
)

// opendict is 국립국어원 우리말샘 (opendict.korean.go.kr). Its autocomplete
// endpoint returns a JSON envelope whose first element is a snippet of
// JavaScript declaring the suggestion array.
const (
	autoCompleteURL = "https://opendict.korean.go.kr/search/autoComplete"
	resultURL       = "https://opendict.korean.go.kr/search/searchResult?focus_name=query&query=%s"
	searchCacheAge  = 30 * time.Second
)

// suggestionsRE extracts the array literal from the response JavaScript:
//
//	var dq_searchResultList=new Array('사랑','사랑가', … );
var suggestionsRE = regexp.MustCompile(`var dq_searchResultList=new Array\((.*?)\);`)

// search is the 우리말샘 autocomplete handler, ported from korean_search.py.
func search(fb *alfred.Feedback, word string) error {
	if word == "" {
		return nil
	}

	// Head item: run the raw query as typed.
	fb.Add(alfred.ItemOpts{
		Title:        "'" + word + "' 검색하기",
		Autocomplete: word,
		Arg:          word,
		QuicklookURL: quick(resultURL, word),
		Valid:        true,
	})

	body, err := cache.Cached(cache.Key("opendict", word), searchCacheAge, func() ([]byte, error) {
		return httpx.Post(autoCompleteURL, map[string]string{"searchTerm": word}, nil)
	})
	if err != nil {
		return err
	}

	// Envelope: {"json": ["var …new Array('a','b'…);", 371]}.
	var res struct {
		JSON []any `json:"json"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	js, _ := firstString(res.JSON)
	suggestions := parseSuggestions(js)
	if len(suggestions) == 0 {
		fb.Add(alfred.ItemOpts{
			Title: "'" + word + "'에 대한 검색 결과가 없습니다",
			Icon:  iconNoResults,
			Valid: false,
		})
		return nil
	}

	for _, s := range suggestions {
		fb.Add(alfred.ItemOpts{
			Title:        s,
			Autocomplete: s,
			Arg:          s,
			Copy:         s,
			LargeType:    s,
			QuicklookURL: quick(resultURL, s),
			Valid:        true,
		})
	}
	return nil
}

// firstString returns the first element of arr coerced to a string.
func firstString(arr []any) (string, bool) {
	if len(arr) == 0 {
		return "", false
	}
	s, ok := arr[0].(string)
	return s, ok
}

// parseSuggestions pulls the suggestion words out of the response JavaScript,
// mirroring parse_suggestions() in korean_search.py: match the Array literal,
// split on commas, and strip the surrounding single quotes. Empty entries
// (e.g. an empty Array()) are dropped.
func parseSuggestions(js string) []string {
	m := suggestionsRE.FindStringSubmatch(js)
	if m == nil {
		return nil
	}
	var out []string
	for _, part := range strings.Split(m[1], ",") {
		s := strings.Trim(strings.TrimSpace(part), "'")
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
