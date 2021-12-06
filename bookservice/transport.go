package Bookservice

//is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a booksvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter() //Ascoltatore delle api
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    //                          adds another Book
	// GET     /books/:id                       retrieves the given Book by id
	// PUT     /books/:id                       post updated Book information about the Book
	// PATCH   /books/:id                       partial updated Book information
	// DELETE  /books/:id                       remove the given Book

	r.Methods("POST").Path("/books/").Handler(httptransport.NewServer(
		e.PostBookEndpoint,
		decodePostBookRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/books/{id}").Handler(httptransport.NewServer(
		e.GetBookEndpoint,
		decodeGetBookRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/books/{id}").Handler(httptransport.NewServer(
		e.PutBookEndpoint,
		decodePutBookRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/books/{id}").Handler(httptransport.NewServer(
		e.PatchBookEndpoint,
		decodePatchBookRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/books/{id}").Handler(httptransport.NewServer(
		e.DeleteBookEndpoint,
		decodeDeleteBookRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postBookRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Book); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getBookRequest{ID: id}, nil
}

func decodePutBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var Book Book
	if err := json.NewDecoder(r.Body).Decode(&Book); err != nil {
		return nil, err
	}
	return putBookRequest{
		ID:   id,
		Book: Book,
	}, nil
}

func decodePatchBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var Book Book
	if err := json.NewDecoder(r.Body).Decode(&Book); err != nil {
		return nil, err
	}
	return patchBookRequest{
		ID:   id,
		Book: Book,
	}, nil
}

func decodeDeleteBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteBookRequest{ID: id}, nil
}

func decodeGetAddressesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getAddressesRequest{BookID: id}, nil
}

func decodeGetAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	addressID, ok := vars["addressID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getAddressRequest{
		BookID:    id,
		AddressID: addressID,
	}, nil
}

func decodePostAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var address Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		return nil, err
	}
	return postAddressRequest{
		BookID:  id,
		Address: address,
	}, nil
}

func decodeDeleteAddressRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	addressID, ok := vars["addressID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteAddressRequest{
		BookID:    id,
		AddressID: addressID,
	}, nil
}

func encodePostBookRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/books/")
	req.URL.Path = "/books/"
	return encodeRequest(ctx, req, request)
}

func encodeGetBookRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/books/{id}")
	r := request.(getBookRequest)
	BookID := url.QueryEscape(r.ID)
	req.URL.Path = "/books/" + BookID
	return encodeRequest(ctx, req, request)
}

func encodePutBookRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/books/{id}")
	r := request.(putBookRequest)
	BookID := url.QueryEscape(r.ID)
	req.URL.Path = "/books/" + BookID
	return encodeRequest(ctx, req, request)
}

func encodePatchBookRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/books/{id}")
	r := request.(patchBookRequest)
	BookID := url.QueryEscape(r.ID)
	req.URL.Path = "/books/" + BookID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteBookRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/books/{id}")
	r := request.(deleteBookRequest)
	BookID := url.QueryEscape(r.ID)
	req.URL.Path = "/books/" + BookID
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressesRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/books/{id}/addresses/")
	r := request.(getAddressesRequest)
	BookID := url.QueryEscape(r.BookID)
	req.URL.Path = "/books/" + BookID + "/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeGetAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/books/{id}/addresses/{addressID}")
	r := request.(getAddressRequest)
	BookID := url.QueryEscape(r.BookID)
	addressID := url.QueryEscape(r.AddressID)
	req.URL.Path = "/books/" + BookID + "/addresses/" + addressID
	return encodeRequest(ctx, req, request)
}

func encodePostAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/books/{id}/addresses/")
	r := request.(postAddressRequest)
	BookID := url.QueryEscape(r.BookID)
	req.URL.Path = "/books/" + BookID + "/addresses/"
	return encodeRequest(ctx, req, request)
}

func encodeDeleteAddressRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/books/{id}/addresses/{addressID}")
	r := request.(deleteAddressRequest)
	BookID := url.QueryEscape(r.BookID)
	addressID := url.QueryEscape(r.AddressID)
	req.URL.Path = "/books/" + BookID + "/addresses/" + addressID
	return encodeRequest(ctx, req, request)
}

func decodePostBookResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postBookResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetBookResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getBookResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutBookResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putBookResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchBookResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response patchBookResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteBookResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteBookResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressesResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getAddressesResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePostAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteAddressResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteAddressResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// booksvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
