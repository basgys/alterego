package router

import (
	"net/url"
	"strings"
)

// Router contains a set of redirections
type Router struct {
	Rules []Rule
}

// New initialises a new router
func New(rules []Rule) *Router {
	return &Router{Rules: rules}
}

// Match tries to find a matching redirection for the URL req
func (rt *Router) Match(req *url.URL) (*url.URL, bool) {
	for _, rule := range rt.Rules {
		if (rule.Src.Host == "" || rule.Src.Host == req.Host) &&
			(rule.Src.Path == "" || rule.Src.Path == req.Path) {
			var path string
			if rule.Dst.Path != "" {
				path = rule.Dst.Path
			} else {
				path = req.Path
			}
			path = strings.TrimSuffix(path, "/")

			return &url.URL{
				Scheme:   rule.Dst.Scheme,
				Host:     rule.Dst.Host,
				Path:     path,
				RawQuery: req.RawQuery,
			}, true
		}
	}

	return nil, false
}

type Rule struct {
	Src url.URL
	Dst url.URL
}
