// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Error defines model for Error.
type Error struct {
	Code   string                  `json:"code"`
	Detail string                  `json:"detail"`
	Meta   *map[string]interface{} `json:"meta,omitempty"`
	Status int                     `json:"status"`
	Title  string                  `json:"title"`
}

// ErrorResponse Response that contains the list of errors
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// ExchangeRateRequestBody defines model for ExchangeRateRequestBody.
type ExchangeRateRequestBody struct {
	CurrencyPair string `json:"currency_pair"`
}

// ExchangeRateResponseData defines model for ExchangeRateResponseData.
type ExchangeRateResponseData map[string]float32

// ExchangeRateResponse defines model for ExchangeRateResponse.
type ExchangeRateResponse struct {
	Data ExchangeRateResponseData `json:"data"`
}

// V1GetExchangeRateJSONBody defines parameters for V1GetExchangeRate.
type V1GetExchangeRateJSONBody struct {
	Data *ExchangeRateRequestBody `json:"data,omitempty"`
}

// V1GetExchangeRateJSONRequestBody defines body for V1GetExchangeRate for application/json ContentType.
type V1GetExchangeRateJSONRequestBody V1GetExchangeRateJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get exchange rate for a currency pair
	// (POST /v1/exchange-rate)
	V1GetExchangeRate(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Get exchange rate for a currency pair
// (POST /v1/exchange-rate)
func (_ Unimplemented) V1GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// V1GetExchangeRate operation middleware
func (siw *ServerInterfaceWrapper) V1GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.V1GetExchangeRate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/exchange-rate", wrapper.V1GetExchangeRate)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xVW2/jOA/9KwK/79GNncwOMOunnUsxKLDYLVp0XxZBwcpMooEteSg6aFDkvy8k2XEu",
	"LjpPfbN14TmHPCJfQLumdZaseChfgMm3znqKP9fPeoN2TXcodNdvhHXtrJCV8IltWxuNYpzNf3hnw5rX",
	"G2owfLXsWmIxKVyFElf/z7SCEv6Xj9B5uuPzKchv4d5+nwHTz84wVVD+m4ItM5BdS1CCe/pBWmAfjlXk",
	"NZs2cIISrp+xaWtSg7IYqYeLIpkdX5LVropa+/he2Ng1xOiCpp7caigJPOOUgReUzh9tGSu0Jg57YqSe",
	"QjrTm44d4A8xs8T0IhMZPF9Rkt7XMqVhPooDT7wlfqSYgVEZ3BNvjSYl1LSOkU29U53FLZoan2rKFJPw",
	"TtUoFO4NsjV2nqrHpx2U8LVG7//Cho7lfyyKg94IQqwSeKxbrMSxz04LOewo2aCo4EE01ivZkKqNF+VW",
	"KVjIyWkt++XyBYxQ49/0YKJ0yCgy4+6iIH3Q5ajolH92aYNTc//syMsXV+0mzNcxk9W7xxZN9GZfSSjh",
	"4f7b1fXD3Rj+Fb+chjhm+QqJN/kePcbw8KvKhMJgfXua64FpMft9kcHKcYMCJaxqhzKC2K55Sv4/XHg5",
	"aIt395Nv29iVG1oQ6tiCqEmu1RvkmryujRVnF0Xx4Y912Jpp18BFW/h8e6NWjtWKRG+MXSsmrK/ENKSo",
	"l60YhXw8NWRThWz6GYzp/DrsDMlSIVvq8+0NZLAl9gluPitmRWDhWrLYGijhw2w+K4JXUTYxdfl2ng/Y",
	"VwE7GsN5mXoLwoa2FN1/wjfSRbU2W7KntAPrUKnYrG8qKOGf+XeS4yJD8tCRL9+t1Y+o++nKj+YW7igu",
	"HI2qRVG8hnc4NzlcQk1+S5d/WeubzWMMfjmPvmClerkRe7F4P+wH27LT5H1o4+raipFdIPHxPRNwY4XY",
	"Yq36CdC32zAnuqZBDuPjO8mkrU8MHV4hrn3odsZundHkYRkB02ALO+cP55Zd1enwo9IhyKDj0D42Iq0v",
	"8xxbM9NYoUVs2753nEe5F1yHnvFKCJ+2Z78S6k+nsVYVbal2bUNWpoKWeV6HcxvnpfxUfCpgvzyIP4/4",
	"9/DGvWIKM7pS4lIvgwxsmMklpN/9cv9fAAAA//8HqQiyAQoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
