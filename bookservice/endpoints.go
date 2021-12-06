package bookservice

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a Book service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	PostBookEndpoint   endpoint.Endpoint
	GetBookEndpoint    endpoint.Endpoint
	PutBookEndpoint    endpoint.Endpoint
	PatchBookEndpoint  endpoint.Endpoint
	DeleteBookEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a Booksvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostBookEndpoint:   MakePostBookEndpoint(s),
		GetBookEndpoint:    MakeGetBookEndpoint(s),
		PutBookEndpoint:    MakePutBookEndpoint(s),
		PatchBookEndpoint:  MakePatchBookEndpoint(s),
		DeleteBookEndpoint: MakeDeleteBookEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a Booksvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		PostBookEndpoint:   httptransport.NewClient("POST", tgt, encodePostBookRequest, decodePostBookResponse, options...).Endpoint(),
		GetBookEndpoint:    httptransport.NewClient("GET", tgt, encodeGetBookRequest, decodeGetBookResponse, options...).Endpoint(),
		PutBookEndpoint:    httptransport.NewClient("PUT", tgt, encodePutBookRequest, decodePutBookResponse, options...).Endpoint(),
		PatchBookEndpoint:  httptransport.NewClient("PATCH", tgt, encodePatchBookRequest, decodePatchBookResponse, options...).Endpoint(),
		DeleteBookEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteBookRequest, decodeDeleteBookResponse, options...).Endpoint(),
	}, nil
}

// PostBook implements Service. Primarily useful in a client.
func (e Endpoints) PostBook(ctx context.Context, b Book) error {
	request := postBookRequest{Book: b}
	response, err := e.PostBookEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postBookResponse)
	return resp.Err
}

// GetBook implements Service. Primarily useful in a client.
func (e Endpoints) GetBook(ctx context.Context, id string) (Book, error) {
	request := getBookRequest{ID: id}
	response, err := e.GetBookEndpoint(ctx, request)
	if err != nil {
		return Book{}, err
	}
	resp := response.(getBookResponse)
	return resp.Book, resp.Err
}

// PutBook implements Service. Primarily useful in a client.
func (e Endpoints) PutBook(ctx context.Context, id string, b Book) error {
	request := putBookRequest{ID: id, Book: b}
	response, err := e.PutBookEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putBookResponse)
	return resp.Err
}

// PatchBook implements Service. Primarily useful in a client.
func (e Endpoints) PatchBook(ctx context.Context, id string, b Book) error {
	request := patchBookRequest{ID: id, Book: b}
	response, err := e.PatchBookEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchBookResponse)
	return resp.Err
}

// DeleteBook implements Service. Primarily useful in a client.
func (e Endpoints) DeleteBook(ctx context.Context, id string) error {
	request := deleteBookRequest{ID: id}
	response, err := e.DeleteBookEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteBookResponse)
	return resp.Err
}



// MakePostBookEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postBookRequest)
		e := s.PostBook(ctx, req.Book)
		return postBookResponse{Err: e}, nil
	}
}

// MakeGetBookEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getBookRequest)
		b, e := s.GetBook(ctx, req.ID)
		return getBookResponse{Book: p, Err: e}, nil
	}
}

// MakePutBookEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putBookRequest)
		e := s.PutBook(ctx, req.ID, req.Book)
		return putBookResponse{Err: e}, nil
	}
}

// MakePatchBookEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchBookRequest)
		e := s.PatchBook(ctx, req.ID, req.Book)
		return patchBookResponse{Err: e}, nil
	}
}

// MakeDeleteBookEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteBookRequest)
		e := s.DeleteBook(ctx, req.ID)
		return deleteBookResponse{Err: e}, nil
	}
}



// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type postBookRequest struct {
	Book Book
}

type postBookResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postBookResponse) error() error { return r.Err }

type getBookRequest struct {
	ID string
}

type getBookResponse struct {
	Book Book `json:"Book,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getBookResponse) error() error { return r.Err }

type putBookRequest struct {
	ID      string
	Book Book
}

type putBookResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putBookResponse) error() error { return nil }

type patchBookRequest struct {
	ID      string
	Book Book
}

type patchBookResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchBookResponse) error() error { return r.Err }

type deleteBookRequest struct {
	ID string
}

type deleteBookResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteBookResponse) error() error { return r.Err }