package core

import (
	"fmt"
	"strings"

	"github.com/awesome-goose/contracts"
	str "github.com/awesome-goose/utils/string"
)

type router struct{}

func NewRouter() *router {
	return &router{}
}

func (r *router) Find(routes contracts.Routes, segments []string) (*contracts.Route, error) {
	var mergedMiddlewares []contracts.Middleware
	var mergedValidators []contracts.Validator
	var current contracts.Route

	for len(segments) > 0 {
		found := false

		method := segments[0]
		isValidMethod := str.IsValidHTTPMethod(method)
		if isValidMethod {
			segments = segments[1:]
		}

		for _, route := range routes {
			if isValidMethod && route.Method == "" {
				continue
			}

			if isValidMethod && route.Method != "" && route.Method != contracts.Method(method) {
				continue
			}

			match, consumed := r.match(route.Path, segments)
			if match {
				current = route
				mergedMiddlewares = append(mergedMiddlewares, route.Middlewares...)
				mergedValidators = append(mergedValidators, route.Validators...)
				segments = segments[consumed:]
				routes = route.Children
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("no route match for segments: %v", segments)
		}
	}

	current.Middlewares = mergedMiddlewares
	current.Validators = mergedValidators
	return &current, nil
}

func (r *router) match(path string, segments []string) (bool, int) {
	routeSegs := r.split(path)
	if len(routeSegs) > len(segments) {
		return false, 0
	}

	for i, seg := range routeSegs {
		if strings.HasPrefix(seg, ":") {
			continue
		}
		if seg != segments[i] {
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
