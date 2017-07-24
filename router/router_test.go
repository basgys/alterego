package router_test

import (
	"net/url"
	"testing"

	"github.com/basgys/alterego/router"
)

func TestRouter_Match(t *testing.T) {
	table := []struct {
		Input       *url.URL
		Rules       []router.Rule
		ExpectURL   *url.URL
		ExpectMatch bool
	}{
		{
			Input: mustParseURL("http://localhost:3000"),
			Rules: []router.Rule{
				router.Rule{
					Src: *mustParseURL("http://localhost:3000"),
					Dst: *mustParseURL("https://stairlin.com"),
				},
			},
			ExpectURL:   mustParseURL("https://stairlin.com"),
			ExpectMatch: true,
		},
		{
			Input: mustParseURL("http://localhost:3000/baidu"),
			Rules: []router.Rule{
				router.Rule{
					Src: *mustParseURL("http://localhost:3000/google"),
					Dst: *mustParseURL("https://google.com/"),
				},
				router.Rule{
					Src: *mustParseURL("http://localhost:3000/baidu"),
					Dst: *mustParseURL("https://baidu.com/"),
				},
			},
			ExpectURL:   mustParseURL("https://baidu.com"),
			ExpectMatch: true,
		},
		{
			Input: mustParseURL("http://localhost:3000/baidu"),
			Rules: []router.Rule{
				router.Rule{
					Src: *mustParseURL("http://localhost:3000/baidu"),
					Dst: *mustParseURL("https://baidu.com/"),
				},
				router.Rule{
					Src: *mustParseURL("http://localhost:3000"),
					Dst: *mustParseURL("https://stairlin.com/"),
				},
			},
			ExpectURL:   mustParseURL("https://baidu.com"),
			ExpectMatch: true,
		},
		{
			Input: mustParseURL("http://localhost:3000/baidu?q=喂"),
			Rules: []router.Rule{
				router.Rule{
					Src: *mustParseURL("http://localhost:3000/google"),
					Dst: *mustParseURL("https://google.com/"),
				},
				router.Rule{
					Src: *mustParseURL("http://localhost:3000/baidu"),
					Dst: *mustParseURL("https://baidu.com/"),
				},
			},
			ExpectURL:   mustParseURL("https://baidu.com?q=喂"),
			ExpectMatch: true,
		},
	}

	for _, test := range table {
		rt := router.New(test.Rules)
		url, ok := rt.Match(test.Input)
		if test.ExpectMatch != ok {
			t.Errorf("expect match %t", test.ExpectMatch)
		}
		if !ok {
			continue
		}
		if test.ExpectURL.String() != url.String() {
			t.Errorf("expect url %s, but got %s", test.ExpectURL.String(), url.String())
		}
	}
}

func BenchmarkRouter_Match(b *testing.B) {
	rt := router.New([]router.Rule{
		router.Rule{
			Src: *mustParseURL("http://localhost:3000/baidu"),
			Dst: *mustParseURL("https://baidu.com/"),
		},
		router.Rule{
			Src: *mustParseURL("http://localhost:3000"),
			Dst: *mustParseURL("https://stairlin.com/"),
		},
		router.Rule{
			Src: *mustParseURL("http://localhost:4000"),
			Dst: *mustParseURL("https://stairlin.com/"),
		},
	})
	requests := []*url.URL{
		mustParseURL("http://localhost:3000/baidu"),
		mustParseURL("http://localhost:3000"),
		mustParseURL("http://localhost:3000/google"),
		mustParseURL("http://localhost:3000/google?q=hello?"),
		mustParseURL("http://stairlin.com"),
	}

	for n := 0; n < b.N; n++ {
		r := requests[n%len(requests)]
		rt.Match(r)
	}
}

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
