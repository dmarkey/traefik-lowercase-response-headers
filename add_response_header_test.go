package traefik_plugin_lowercase_response_headers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

type dummyHandler struct{}

func (dummyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("overwrite", "true")
	w.WriteHeader(200)
}

func (s *Suite) TestMain() {
	cfg := CreateConfig()

	data := "123bla321"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	resp.Header().Set("FOO", data)
	h.ServeHTTP(resp, req)
	s.Require().Equal(data, resp.Header().Get("foo"))
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
