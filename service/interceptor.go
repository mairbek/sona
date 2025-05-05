package service

import (
	"fmt"
	"net/http"
	"strings"
)

// PathRewriteRule defines a single path rewrite rule
type PathRewriteRule struct {
	From string
	To   string
}

// PathRewriteInterceptor is a middleware that rewrites request paths based on configured rules
type PathRewriteInterceptor struct {
	rules []PathRewriteRule
	next  http.Handler
}

// NewPathRewriteInterceptor creates a new PathRewriteInterceptor with the given rules
func NewPathRewriteInterceptor(rules []PathRewriteRule, next http.Handler) *PathRewriteInterceptor {
	return &PathRewriteInterceptor{
		rules: rules,
		next:  next,
	}
}

// ServeHTTP implements http.Handler interface
func (i *PathRewriteInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check each rule to see if it matches the current path
	fmt.Printf("--- fml %v", r.URL.Path)
	for _, rule := range i.rules {
		if strings.HasPrefix(r.URL.Path, rule.From) {
			// Rewrite the path according to the rule
			r.URL.Path = rule.To
			break
		}
	}

	// Pass the request to the next handler
	i.next.ServeHTTP(w, r)
}
