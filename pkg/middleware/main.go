package middleware

import "net/http"

type Middleware func(h http.HandlerFunc) http.HandlerFunc


func MiddlewareIn(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >=0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func MiddlewareOut(h http.HandlerFunc, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
}
//Middleware Chains


