package server

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerCondition func(*http.Request) bool

type scPair struct {
	server    http.Handler
	condition ServerCondition
}

type ServerMux struct {
	pairs []*scPair
}

func (sm *ServerMux) Add(server http.Handler, condition ServerCondition) {
	if server == nil {
		return
	}
	sm.pairs = append(sm.pairs, &scPair{server: server, condition: condition})
}

func (sm *ServerMux) Remove(server http.Handler) {
	if server == nil {
		return
	}
	for i, pair := range sm.pairs {
		if pair.server == server {
			sm.pairs = append(sm.pairs[:i], sm.pairs[i+1:]...)
		}
	}
}

func (sm *ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, pair := range sm.pairs {
		s, cond := pair.server, pair.condition
		if cond != nil && !cond(r) {
			continue
		}
		s.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func NewServerMux(server ...any) *ServerMux {
	if len(server)%2 != 0 {
		panic(fmt.Sprintf("server: len of params passed must be a multiple of 2, got %d", len(server)))
	}

	sm := new(ServerMux)
	for i := 0; i < len(server); i += 2 {
		var h http.Handler
		var cond ServerCondition
		var ok bool
		if h, ok = server[i].(http.Handler); !ok {
			panic(fmt.Sprintf("server: type of server isn't http.Handler, got %T", h))
		}
		if cond, ok = server[i+1].(func(*http.Request) bool); server[i+1] != nil && !ok {
			panic(fmt.Sprintf("server: type of condition isn't func(*http.Request) bool, got %T", cond))
		}
		sm.Add(h, cond)
	}
	return sm
}

func EnableH2C(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}
