package http

import (
	"context"
	"encoding/json"
	"github.com/stickysh/sticky/pkg/endpoint"
	"net/http"
)

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error
type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)
type ErrorEncFunc func(context.Context, error, http.ResponseWriter)

type server struct {
	e endpoint.Endpoint
	dec DecodeRequestFunc
	enc EncodeResponseFunc
	errEnc ErrorEncFunc
}

func NewServer(e endpoint.Endpoint, dec DecodeRequestFunc, enc EncodeResponseFunc) *server {
	return &server{
		e: e,
		dec: dec,
		enc: enc,
		errEnc: ErrorEncoder,
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec, err := s.dec(ctx, r)
	if err != nil {
		s.errEnc(ctx, err, w)
		return
	}

	result, err := s.e(ctx, dec)
	if err != nil {
		s.errEnc(ctx, err, w)
		return
	}

	if err := s.enc(ctx, w, result); err != nil {
		s.errEnc(ctx, err, w)
		return
	}

}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := http.StatusInternalServerError

	var body struct {
		Error string `json:"error,omitempty"`
	}
	body.Error = err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}