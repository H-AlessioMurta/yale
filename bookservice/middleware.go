package bookservices

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostBook(ctx context.Context, b Book) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostBook", "id", b.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostBook(ctx, b)
}

func (mw loggingMiddleware) GetBook(ctx context.Context, id string) (b Book, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetBook", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetBook(ctx, id)
}

func (mw loggingMiddleware) PutBook(ctx context.Context, id string, b Book) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutBook", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutBook(ctx, id, b)
}

func (mw loggingMiddleware) PatchBook(ctx context.Context, id string, b Book) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchBook", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchBook(ctx, id, b)
}

func (mw loggingMiddleware) DeleteBook(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteBook", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteBook(ctx, id)
}
