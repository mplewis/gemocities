package router

import (
	"context"
	"net/url"

	"git.sr.ht/~adnano/go-gemini"
)

type Params = map[string]string
type RouteFunction func(ctx context.Context, w gemini.ResponseWriter, rq Request)

type Router struct {
	routes []Route
}

type RouterMatch struct {
	Router Router
	Params Params
}

type Request struct {
	Raw         *gemini.Request
	PathParams  Params
	QueryParams Params
}

func NewRouter(routes ...Route) Router {
	return Router{routes}
}

func (r Router) ServeGemini(ctx context.Context, w gemini.ResponseWriter, rq *gemini.Request) {
	var found bool
	var route Route
	var pathParams Params
	for _, cand := range r.routes {
		pp, match := cand.Match(rq.URL.Path)
		if match {
			found = true
			route = cand
			pathParams = pp
			break
		}
	}

	if !found {
		w.WriteHeader(gemini.StatusNotFound, "path not found")
		return
	}

	route.handler(ctx, w, Request{
		PathParams:  pathParams,
		QueryParams: flattenQueryParams(rq.URL.Query()),
		Raw:         rq,
	})
}

func flattenQueryParams(raw url.Values) Params {
	params := Params{}
	for k, v := range raw {
		params[k] = v[0]
	}
	return params
}
