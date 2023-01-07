package traefik_plugin_add_response_header

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

type Config struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type plugin struct {
	name   string
	next   http.Handler
	config *Config
}

func (p *plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := httptest.NewRecorder()

	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: w headers before - %+v", w.Header()))

	p.next.ServeHTTP(resp, req)

	for k := range resp.Header() {
		w.Header().Set(k, resp.Header().Get(k))
	}
	w.Header().Set(p.config.To, req.Header.Get(p.config.From))

	io.Copy(w, resp.Body)

	w.WriteHeader(resp.Code)
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.From == "" {
		return nil, fmt.Errorf("from cannot be empty")
	}
	if config.To == "" {
		return nil, fmt.Errorf("to cannot be empty")
	}

	return &plugin{
		name:   name,
		next:   next,
		config: config,
	}, nil
}
