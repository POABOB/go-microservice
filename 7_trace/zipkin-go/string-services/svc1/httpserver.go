//go:build go1.7
// +build go1.7

package svc1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"

	"zipkin-go/middleware"
)

type httpService struct {
	service Service
}

// 處理 Concat 的參數請求
func (s *httpService) concatHandler(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	result, err := s.service.Concat(req.Context(), v.Get("a"), v.Get("b"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(result))
}

// 處理 Sum 的參數請求
func (s *httpService) sumHandler(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	a, err := strconv.ParseInt(v.Get("a"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := strconv.ParseInt(v.Get("b"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.service.Sum(req.Context(), a, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d", result)))
}

// NewHTTPHandler returns a new HTTP handler our svc1.
func NewHTTPHandler(tracer opentracing.Tracer, service Service) http.Handler {
	// Create our HTTP Service.
	svc := &httpService{service: service}

	// Create the mux.
	mux := http.NewServeMux()

	// Concat Handler
	var concatHandler http.Handler = http.HandlerFunc(svc.concatHandler)
	// 使用 tracing middleware 包裝 Concat Handler
	concatHandler = middleware.FromHTTPRequest(tracer, "Concat")(concatHandler)

	// Sum handler.
	var sumHandler http.Handler = http.HandlerFunc(svc.sumHandler)
	// 使用 tracing middleware 包裝 Sum Handler
	sumHandler = middleware.FromHTTPRequest(tracer, "Sum")(sumHandler)

	// Wire up the mux.
	mux.Handle("/concat/", concatHandler)
	mux.Handle("/sum/", sumHandler)
	return mux
}
