package core

import (
	"strings"

	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/errors"
	str "github.com/awesome-goose/utils/string"
)

type router struct{}

func NewRouter() *router {
	return &router{}
}

func (r *router) Find(routes contracts.Routes, method string, paths []string) (*contracts.Route, error) {
	var mergedMiddlewares []contracts.Middleware
	var mergedValidators []contracts.Validator
	var current contracts.Route

	for len(paths) > 0 {
		found := false

		isValidMethod := str.IsValidHTTPMethod(method)
		for _, route := range routes {
			if isValidMethod && route.Method == "" {
				continue
			}

			if isValidMethod && route.Method != "" && route.Method != contracts.Method(method) {
				continue
			}

			match, consumed := r.match(route.Path, paths)
			if match {
				current = route
				mergedMiddlewares = append(mergedMiddlewares, route.Middlewares...)
				mergedValidators = append(mergedValidators, route.Validators...)
				paths = paths[consumed:]
				routes = route.Children
				found = true
				break
			}
		}

		if !found {
			return nil, errors.ErrNoRouteMatch.WithMeta(paths)
		}
	}

	current.Middlewares = mergedMiddlewares
	current.Validators = mergedValidators
	return &current, nil
}

func (r *router) match(path string, paths []string) (bool, int) {
	routeSegs := r.split(path)
	if len(routeSegs) > len(paths) {
		return false, 0
	}

	for i, seg := range routeSegs {
		if strings.HasPrefix(seg, ":") {
			continue
		}
		if seg != paths[i] {
			return false, 0
		}
	}
	return true, len(routeSegs)
}

func (r *router) split(path string) []string {
	raw := strings.Split(path, "/")
	var cleaned []string
	for _, seg := range raw {
		seg = strings.TrimSpace(seg)
		if seg != "" {
			cleaned = append(cleaned, seg)
		}
	}
	return cleaned
}
