package traefik_plugin_lowercase_response_headers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

var (
	_ interface {
		http.ResponseWriter
		http.Hijacker
	} = &wrappedResponseWriter{}
)

type Config struct {
}

func CreateConfig() *Config {
	return &Config{}
}

type plugin struct {
	name   string
	next   http.Handler
	config *Config
	regex  *regexp.Regexp
}

type wrappedResponseWriter struct {
	w    http.ResponseWriter
	buf  *bytes.Buffer
	code int
}

func (w *wrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.code = code
}

func (w *wrappedResponseWriter) Flush() {
	w.w.WriteHeader(w.code)
	io.Copy(w.w, w.buf)
}

func (w *wrappedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not an http.Hijacker", w.w)
	}

	return hijacker.Hijack()
}

func (p *plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := &wrappedResponseWriter{
		w:    w,
		buf:  &bytes.Buffer{},
		code: 200,
	}
	defer resp.Flush()

	p.next.ServeHTTP(resp, req)

	for name, values := range resp.Header() {
		resp.Header().Del(name)
		resp.Header().Add(strings.ToLower(name), values[0])
	}

}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &plugin{
		name:   name,
		next:   next,
		config: config,
	}, nil
}
