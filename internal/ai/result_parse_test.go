package ai

import "testing"

func TestParseResultContent(t *testing.T) {
	cases := []struct {
		name    string
		content string
		ok      bool
		title   string
	}{
		{
			name:    "plain json",
			content: `{"title":"A","summary":"B","sections":[{"title":"S","items":["i"]}],"warnings":["w"],"raw":""}`,
			ok:      true,
			title:   "A",
		},
		{
			name: "json in markdown fence",
			content: "```json\n" +
				`{"title":"Fence","summary":"B","sections":[{"title":"S","items":["i"]}],"warnings":["w"],"raw":""}` +
				"\n```",
			ok:    true,
			title: "Fence",
		},
		{
			name: "json with leading text",
			content: "Tady je výstup:\n\n" +
				`{"title":"Lead","summary":"B","sections":[{"title":"S","items":["i"]}],"warnings":["w"],"raw":""}`,
			ok:    true,
			title: "Lead",
		},
		{
			name:    "non json",
			content: "ahoj světe",
			ok:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := parseResultContent(tc.content)
			if ok != tc.ok {
				t.Fatalf("expected ok=%v, got %v", tc.ok, ok)
			}
			if tc.ok && result.Title != tc.title {
				t.Fatalf("expected title %q, got %q", tc.title, result.Title)
			}
		})
	}
}
