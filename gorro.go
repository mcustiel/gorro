package gorro

import (
	"net/http"
)

type Request struct {
	http.Request
	NamedParams map[string]string
	Params      []string
	RegexGroups []string
	Regex       string
}

type Handler func(w http.ResponseWriter, r *Request) error

type HandlersMap map[string]Handler;

type Router interface {
	Register(regex string, handlers HandlersMap) error
	Route(w http.ResponseWriter, r *http.Request) error
	OnNotFound(handler func(w http.ResponseWriter, r *http.Request))
	OnError(handler func(w http.ResponseWriter, r *http.Request, err error))
}
