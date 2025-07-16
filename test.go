package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vanern/goapi/typs"
)

// AuthMiddleware returns a Middleware that enforces:
//   - roles == nil               → public (no auth)
//   - roles == []string{AuthOnly}→ valid token only
//   - len(roles)>0               → valid token + at least one required role
func Test1() typs.Middleware {
	return func(ep typs.Endpoint, next http.Handler) http.Handler {
		return mw1(ep, next)
	}
}

func Test2() typs.Middleware {
	return func(ep typs.Endpoint, next http.Handler) http.Handler {
		return mw2(ep, next)
	}
}

func mw1(ep typs.Endpoint, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Applying Test middleware 1 for endpoint:", ep.Method, ep.Path)
		next.ServeHTTP(w, r)
	})
}

func mw2(ep typs.Endpoint, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Applying Test middleware 2 for endpoint:", ep.Method, ep.Path)
		next.ServeHTTP(w, r)
	})
}
