package servicequeue

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/go-softwarelab/common/pkg/collection"
	"github.com/go-softwarelab/common/pkg/is"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seq2"
	"github.com/go-softwarelab/common/pkg/to"
)

var ErrEmptyResult = errors.New("service returns an empty result")

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
		return to.ZeroValue[R](), errors.New("no services registered")
	}

	results := seq.MapOrErr(seq.FromSlice(services), func(s S) (collection.Tuple2[string, R], error) {
		res, err := callService(s)
		if err != nil {
			err = fmt.Errorf("service %s responds with error: %w", s.Name(), err)
		}
		return collection.NewTuple2(s.Name(), res), err
	})

	results = takeUntilHaveResult(results)

	results = seq2.Each(results, func(serviceResult collection.Tuple2[string, R], err error) {
		if err != nil {
			logger.Warn("error when calling service",
				slog.String("service.name", serviceResult.A),
				logging.Error(err),
			)
		}
	})

	result := seq2.Reduce(results,
		func(res collection.Tuple2[R, error], serviceResult collection.Tuple2[string, R], err error) collection.Tuple2[R, error] {
			ret := collection.NewTuple2[R, error](serviceResult.B, errors.Join(res.B, err))
			return ret
		},
		collection.Tuple2[R, error]{},
	)

	if is.NotNil(result.A) {
		return result.A, nil
	}
	return result.A, fmt.Errorf("all services failed: %w", result.B)
}

func takeUntilHaveResult[R any](seq iter.Seq2[collection.Tuple2[string, R], error]) iter.Seq2[collection.Tuple2[string, R], error] {
	return func(yield func(collection.Tuple2[string, R], error) bool) {
		for result, err := range seq {
			if err != nil {
				if !yield(result, err) {
					break
				}
				continue
			}
			if is.Nil(result.B) {
				if !yield(result, ErrEmptyResult) {
					break
				}
				continue
			}

			yield(result, nil)
			break
		}
	}
}
