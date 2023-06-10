// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: ab5a6c03d7
// Version Date: 2023-06-04T17:09:20Z

package transport

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// This service
	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	endpts "github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"
)

const contentType = "application/json; charset=utf-8"

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewUserClient
	_ = io.Copy
	_ = errors.Wrap
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(_ context.Context, endpoints endpts.Endpoints, zipkinTracer *gozipkin.Tracer, logger *zap.Logger, loggerLevel zapcore.Level) http.Handler {
	responseEncoder := EncodeHTTPGenericResponse

	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitzap.NewZapSugarLogger(logger, loggerLevel))),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
		zipkinServer,
	}
	m := mux.NewRouter()

	m.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.LoginEndpoint,
		DecodeHTTPLoginZeroRequest,
		responseEncoder,
		serverOptions...,
	))

	m.Methods("POST").Path("/register").Handler(httptransport.NewServer(
		endpoints.RegisterEndpoint,
		DecodeHTTPRegisterZeroRequest,
		responseEncoder,
		serverOptions...,
	))

	m.Methods("GET").Path("/loginWithGoogle").Handler(httptransport.NewServer(
		endpoints.LoginWithGoogleEndpoint,
		DecodeHTTPLoginWithGoogleZeroRequest,
		responseEncoder,
		serverOptions...,
	))

	m.Methods("GET").Path("/loginWithGoogleCallback").Handler(httptransport.NewServer(
		endpoints.LoginWithGoogleCallbackEndpoint,
		DecodeHTTPLoginWithGoogleCallbackZeroRequest,
		responseEncoder,
		serverOptions...,
	))

	m.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoints.HealthCheckEndpoint,
		DecodeHTTPHealthCheckZeroRequest,
		responseEncoder,
		serverOptions...,
	))

	m.Path("/metrics").Handler(promhttp.Handler())

	return m
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(errorWrapper{Error: err.Error()})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// httpError satisfies the Headerer and StatusCoder interfaces in
// package github.com/go-kit/kit/transport/http.
type httpError struct {
	error
	statusCode int
	headers    map[string][]string
}

func (h httpError) StatusCode() int {
	return h.statusCode
}

func (h httpError) Headers() http.Header {
	return h.headers
}

// Server Decode

// DecodeHTTPLoginZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded login request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPLoginZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserLoginRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := encodePathParams(mux.Vars(r))
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if UsernameLoginStrArr, ok := queryParams["username"]; ok {
		UsernameLoginStr := UsernameLoginStrArr[0]
		UsernameLogin := UsernameLoginStr
		req.Username = UsernameLogin
	}

	if PasswordLoginStrArr, ok := queryParams["password"]; ok {
		PasswordLoginStr := PasswordLoginStrArr[0]
		PasswordLogin := PasswordLoginStr
		req.Password = PasswordLogin
	}

	return &req, err
}

// DecodeHTTPRegisterZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded register request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPRegisterZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserRegisterRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := encodePathParams(mux.Vars(r))
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if UsernameRegisterStrArr, ok := queryParams["username"]; ok {
		UsernameRegisterStr := UsernameRegisterStrArr[0]
		UsernameRegister := UsernameRegisterStr
		req.Username = UsernameRegister
	}

	if EmailRegisterStrArr, ok := queryParams["email"]; ok {
		EmailRegisterStr := EmailRegisterStrArr[0]
		EmailRegister := EmailRegisterStr
		req.Email = EmailRegister
	}

	if PasswordRegisterStrArr, ok := queryParams["password"]; ok {
		PasswordRegisterStr := PasswordRegisterStrArr[0]
		PasswordRegister := PasswordRegisterStr
		req.Password = PasswordRegister
	}

	if PassconfRegisterStrArr, ok := queryParams["passconf"]; ok {
		PassconfRegisterStr := PassconfRegisterStrArr[0]
		PassconfRegister := PassconfRegisterStr
		req.Passconf = PassconfRegister
	}

	if BirthdayRegisterStrArr, ok := queryParams["birthday"]; ok {
		BirthdayRegisterStr := BirthdayRegisterStrArr[0]
		BirthdayRegister := BirthdayRegisterStr
		req.Birthday = BirthdayRegister
	}

	if SexRegisterStrArr, ok := queryParams["sex"]; ok {
		SexRegisterStr := SexRegisterStrArr[0]
		SexRegister := SexRegisterStr
		req.Sex = SexRegister
	}

	if PreferedLocationRegisterStrArr, ok := queryParams["preferedLocation"]; ok {
		PreferedLocationRegisterStr := PreferedLocationRegisterStrArr[0]

		var PreferedLocationRegister []string
		if len(PreferedLocationRegisterStrArr) > 1 {
			PreferedLocationRegister = PreferedLocationRegisterStrArr
		} else {
			PreferedLocationRegister = strings.Split(PreferedLocationRegisterStr, ",")
		}
		req.PreferedLocation = PreferedLocationRegister
	}

	return &req, err
}

// DecodeHTTPLoginWithGoogleZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded loginwithgoogle request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPLoginWithGoogleZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserLoginRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := encodePathParams(mux.Vars(r))
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if UsernameLoginWithGoogleStrArr, ok := queryParams["username"]; ok {
		UsernameLoginWithGoogleStr := UsernameLoginWithGoogleStrArr[0]
		UsernameLoginWithGoogle := UsernameLoginWithGoogleStr
		req.Username = UsernameLoginWithGoogle
	}

	if PasswordLoginWithGoogleStrArr, ok := queryParams["password"]; ok {
		PasswordLoginWithGoogleStr := PasswordLoginWithGoogleStrArr[0]
		PasswordLoginWithGoogle := PasswordLoginWithGoogleStr
		req.Password = PasswordLoginWithGoogle
	}

	return &req, err
}

// DecodeHTTPLoginWithGoogleCallbackZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded loginwithgooglecallback request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPLoginWithGoogleCallbackZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserLoginRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := encodePathParams(mux.Vars(r))
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if UsernameLoginWithGoogleCallbackStrArr, ok := queryParams["username"]; ok {
		UsernameLoginWithGoogleCallbackStr := UsernameLoginWithGoogleCallbackStrArr[0]
		UsernameLoginWithGoogleCallback := UsernameLoginWithGoogleCallbackStr
		req.Username = UsernameLoginWithGoogleCallback
	}

	if PasswordLoginWithGoogleCallbackStrArr, ok := queryParams["password"]; ok {
		PasswordLoginWithGoogleCallbackStr := PasswordLoginWithGoogleCallbackStrArr[0]
		PasswordLoginWithGoogleCallback := PasswordLoginWithGoogleCallbackStr
		req.Password = PasswordLoginWithGoogleCallback
	}

	return &req, err
}

// DecodeHTTPHealthCheckZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded healthcheck request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPHealthCheckZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.HealthCheckRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := encodePathParams(mux.Vars(r))
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	marshaller := jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}

	return marshaller.Marshal(w, response.(proto.Message))
}

// Helper functions

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// Tune specific change.
	// also add the request url
	ctx = context.WithValue(ctx, "request-url", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}

// encodePathParams encodes `mux.Vars()` with dot notations into JSON objects
// to be unmarshaled into non-basetype fields.
// e.g. {"book.name": "books/1"} -> {"book": {"name": "books/1"}}
func encodePathParams(vars map[string]string) map[string]string {
	var recur func(path, value string, data map[string]interface{})
	recur = func(path, value string, data map[string]interface{}) {
		parts := strings.SplitN(path, ".", 2)
		key := parts[0]
		if len(parts) == 1 {
			data[key] = value
		} else {
			if _, ok := data[key]; !ok {
				data[key] = make(map[string]interface{})
			}
			recur(parts[1], value, data[key].(map[string]interface{}))
		}
	}

	data := make(map[string]interface{})
	for key, val := range vars {
		recur(key, val, data)
	}

	ret := make(map[string]string)
	for key, val := range data {
		switch val := val.(type) {
		case string:
			ret[key] = val
		case map[string]interface{}:
			m, _ := json.Marshal(val)
			ret[key] = string(m)
		}
	}
	return ret
}