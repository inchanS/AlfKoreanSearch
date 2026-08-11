package handlers

import (
	"reflect"
	"testing"
)

func TestParseSuggestions(t *testing.T) {
	cases := []struct {
		name string
		js   string
		want []string
	}{
		{
			// The case from the original test_workflow.py.
			name: "two words",
			js:   "var dq_searchKeyword='사랑'; var dq_searchResultList=new Array('사랑','사랑하다');",
			want: []string{"사랑", "사랑하다"},
		},
		{
			name: "spaces around entries",
			js:   "var dq_searchResultList=new Array( '가', '나' );",
			want: []string{"가", "나"},
		},
		{
			name: "empty array yields nothing",
			js:   "var dq_searchResultList=new Array();",
			want: nil,
		},
		{
			name: "no match yields nothing",
			js:   "something completely different",
			want: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseSuggestions(c.js); !reflect.DeepEqual(got, c.want) {
				t.Errorf("parseSuggestions() = %#v, want %#v", got, c.want)
			}
		})
	}
}

func TestParseSuggestionsContainsWord(t *testing.T) {
	// Mirror the original unittest assertion: '사랑' is among the results.
	got := parseSuggestions("var dq_searchResultList=new Array('사랑','사랑가','사랑간');")
	found := false
	for _, s := range got {
		if s == "사랑" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected '사랑' in %#v", got)
	}
}

func TestFirstString(t *testing.T) {
	if s, ok := firstString([]any{"js", float64(371)}); !ok || s != "js" {
		t.Errorf("firstString = %q ok=%v", s, ok)
	}
	if _, ok := firstString(nil); ok {
		t.Error("firstString(nil) should not be ok")
	}
	if _, ok := firstString([]any{float64(1)}); ok {
		t.Error("firstString of non-string head should not be ok")
	}
}
