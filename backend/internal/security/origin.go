package security

import (
	"net/http"
	"strings"
)

type OriginChecker struct {
	allowed map[string]struct{}
}

func NewOriginChecker(origins []string) *OriginChecker {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return &OriginChecker{allowed: allowed}
}

func (c *OriginChecker) Check(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return false
	}

	_, ok := c.allowed[origin]
	return ok
}
