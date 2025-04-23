package servicequeue

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/go-softwarelab/common/pkg/is"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/to"
)

var ErrEmptyResult = fmt.Errorf("service returns an empty result")
var ErrNoServicesRegistered = fmt.Errorf("no services registered")

// Queue is a structure that holds a collection of services and abstracts away the details of calling them and error handling.
// Services are functions accepting a context and returning a result or an error.
// For more arguments see Queue1, Queue2, and Queue3.
type Queue[R any] struct {
	logger     *slog.Logger
	methodName string
	services   []*Service[R]
}

func NewQueue[R any](logger *slog.Logger, methodName string, services ...*Service[R]) Queue[R] {
	logger = logging.Child(logger, "services."+methodName)

	return Queue[R]{
		logger:     logger,
		methodName: methodName,
		services:   services,
	}
}

// OneByOne calls the services with provided context, one by one, until a successful result is obtained.
// Returns the first successful result or an error if all services fail.
func (q *Queue[R]) OneByOne(ctx context.Context) (R, error) {
	return processOneByOne(q.logger, q.services, func(s *Service[R]) (R, error) {
		return s.service(ctx)
	})
}

// Queue1 is a structure that holds a collection of services and abstracts away the details of calling them and error handling.
// Services are functions accepting a context and an argument and returning a result or an error.
// For different number of arguments see Queue, Queue2, and Queue3.
type Queue1[A, R any] struct {
	logger     *slog.Logger
	methodName string
	services   []*Service1[A, R]
}

func NewQueue1[A, R any](logger *slog.Logger, methodName string, services ...*Service1[A, R]) Queue1[A, R] {
	logger = logging.Child(logger, "services."+methodName)

	return Queue1[A, R]{
		logger:     logger,
		methodName: methodName,
		services:   services,
	}
}

// OneByOne processes services one by one until a successful result is obtained.
// The context and argument is passed to each service.
// Returns the first successful result or an error if all services fail.
func (q *Queue1[A, R]) OneByOne(ctx context.Context, a A) (R, error) {
	return processOneByOne(q.logger, q.services, func(s *Service1[A, R]) (R, error) {
		return s.service(ctx, a)
	})
}

// Queue2 is a structure that holds a collection of services and abstracts away the details of calling them and error handling.
// Services are functions accepting a context and two arguments and returning a result or an error.
// For different number of arguments see Queue, Queue1, and Queue3.
type Queue2[A, B, R any] struct {
	logger     *slog.Logger
	methodName string
	services   []*Service2[A, B, R]
}

func NewQueue2[A, B, R any](logger *slog.Logger, methodName string, services ...*Service2[A, B, R]) Queue2[A, B, R] {
	logger = logging.Child(logger, "services."+methodName)

	return Queue2[A, B, R]{
		logger:     logger,
		methodName: methodName,
		services:   services,
	}
}

// OneByOne processes services one by one until a successful result is obtained.
// The context and arguments are passed to each service.
// Returns the first successful result or an error if all services fail.
func (q *Queue2[A, B, R]) OneByOne(ctx context.Context, a A, b B) (R, error) {
	return processOneByOne(q.logger, q.services, func(s *Service2[A, B, R]) (R, error) {
		return s.service(ctx, a, b)
	})
}

// Queue3 is a structure that holds a collection of services and abstracts away the details of calling them and error handling.
// Services are functions accepting a context and three arguments and returning a result or an error.
// For different number of arguments see Queue, Queue1, and Queue2.
type Queue3[A, B, C, R any] struct {
	logger     *slog.Logger
	methodName string
	services   []*Service3[A, B, C, R]
}

func NewQueue3[A, B, C, R any](logger *slog.Logger, methodName string, services ...*Service3[A, B, C, R]) Queue3[A, B, C, R] {
	logger = logging.Child(logger, "services."+methodName)

	return Queue3[A, B, C, R]{
		logger:     logger,
		methodName: methodName,
		services:   services,
	}
}

// OneByOne processes services one by one until a successful result is obtained.
// The context and arguments are passed to each service.
// Returns the first successful result or an error if all services fail.
func (q *Queue3[A, B, C, R]) OneByOne(ctx context.Context, a A, b B, c C) (R, error) {
	return processOneByOne(q.logger, q.services, func(s *Service3[A, B, C, R]) (R, error) {
		return s.service(ctx, a, b, c)
	})
}

type serv interface {
	Name() string
}

func processOneByOne[S serv, R any](logger *slog.Logger, services []S, callService func(S) (R, error)) (R, error) {
	if len(services) == 0 {
		return to.ZeroValue[R](), ErrNoServicesRegistered
	}

	results := seq.Map(seq.FromSlice(services), func(s S) serviceCallResult[R] {
		res, err := callService(s)
		return serviceCallResult[R]{
			ServiceName: s.Name(),
			Result:      res,
			Err:         err,
		}
	})

	results = takeUntilHaveResult(results)

	results = seq.Each(results, func(serviceResult serviceCallResult[R]) {
		if serviceResult.Err != nil {
			logger.Warn("error when calling service",
				slog.String("service.name", serviceResult.ServiceName),
				logging.Error(serviceResult.Err),
			)
		}
	})

	var err error
	for result := range results {
		if result.Err != nil {
			err = errors.Join(err, fmt.Errorf("error from service %s: %w", result.ServiceName, result.Err))
			continue
		}
		return result.Result, nil
	}

	return to.ZeroValue[R](), fmt.Errorf("all services failed: %w", err)
}

type serviceCallResult[R any] struct {
	ServiceName string
	Result      R
	Err         error
}

func takeUntilHaveResult[R any](seq iter.Seq[serviceCallResult[R]]) iter.Seq[serviceCallResult[R]] {
	return func(yield func(serviceCallResult[R]) bool) {
		for result := range seq {
			if result.Err != nil {
				errResult := serviceCallResult[R]{
					ServiceName: result.ServiceName,
					Err:         result.Err,
				}

				if !yield(errResult) {
					break
				}
				continue
			}
			if is.Nil(result.Result) {
				nilResult := serviceCallResult[R]{
					ServiceName: result.ServiceName,
					Err:         ErrEmptyResult,
				}

				if !yield(nilResult) {
					break
				}
				continue
			}

			yield(result)
			break
		}
	}
}
