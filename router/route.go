package router

import (
	"fmt"
	"strings"
)

type Route struct {
	parts       []RoutePart
	handler     RouteFunction
	endsInSplat bool
}

type RoutePart struct {
	Key     string
	Matcher bool
	Splat   bool
}

func NewRoute(pattern string, fn RouteFunction) Route {
	raws := strings.Split(strings.Trim(pattern, "/"), "/")
	parts := []RoutePart{}
	for _, raw := range raws {
		if strings.HasPrefix(raw, ":") {
			parts = append(parts, RoutePart{Key: raw[1:], Matcher: true})
		} else if strings.HasPrefix(raw, "*") {
			parts = append(parts, RoutePart{Key: raw[1:], Splat: true})
		} else {
			parts = append(parts, RoutePart{Key: raw})
		}
	}

	endsInSplat := false
	for i, part := range parts {
		if part.Splat {
			if i < len(parts)-1 {
				panic(fmt.Sprintf("bad pattern: `%s`: splat parts must be the last parts in the pattern", pattern))
			}
			endsInSplat = true
		}
	}

	return Route{parts: parts, handler: fn, endsInSplat: endsInSplat}
}

func (r Route) Match(path string) (map[string]string, bool) {
	raws := strings.Split(strings.TrimLeft(path, "/"), "/")
	if r.endsInSplat {
		if len(raws) < len(r.parts) {
			return nil, false
		}
	} else {
		if len(raws) != len(r.parts) {
			return nil, false
		}
	}

	params := map[string]string{}
	for i, a := range raws {
		b := r.parts[i]
		if b.Splat {
			params[b.Key] = "/" + strings.Join(raws[i:], "/")
			break
		} else if b.Matcher {
			params[b.Key] = a
		} else if a != b.Key {
			return nil, false
		}
	}
	return params, true
}
