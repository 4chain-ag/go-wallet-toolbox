package servicequeue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/servicequeue"
	"github.com/go-softwarelab/common/pkg/slices"
	"github.com/stretchr/testify/assert"
)

const secondArgument = "test"
const thirdArgument = 1
const fourthArgument = true

func TestQueue(t *testing.T) {
	tests := map[string]struct {
		services         []TestService
		expectedResult   *Result
		errorExpectation func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool
	}{
		"single service returning success should return success": {
			services: []TestService{
				TestService{Name: "successful"}.Successful(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"single service returning error should return error": {
			services: []TestService{
				TestService{Name: "failing"}.Failing(),
			},
			expectedResult:   nil,
			errorExpectation: assert.Error,
		},
		"single service returning nil result should return error": {
			services: []TestService{
				TestService{Name: "nil-result"}.ReturningNilResult(),
			},
			expectedResult:   nil,
			errorExpectation: assert.Error,
		},
		"all services returning error should return error": {
			services: []TestService{
				TestService{Name: "error-service-1"}.Failing(),
				TestService{Name: "error-service-2"}.Failing(),
				TestService{Name: "error-service-3"}.Failing(),
			},
			expectedResult:   nil,
			errorExpectation: assert.Error,
		},
		"all services returning nil result should return error": {
			services: []TestService{
				TestService{Name: "nil-service-1"}.ReturningNilResult(),
				TestService{Name: "nil-service-2"}.ReturningNilResult(),
				TestService{Name: "nil-service-3"}.ReturningNilResult(),
			},
			expectedResult:   nil,
			errorExpectation: assert.Error,
		},
		"first service returning success should return that result and not call others": {
			services: []TestService{
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called-1"}.ShouldNotBeCalled(),
				TestService{Name: "should-not-be-called-2"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first service returning error but second returning success should return second result": {
			services: []TestService{
				TestService{Name: "error-service"}.Failing(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first service returning nil but second returning success should return second result": {
			services: []TestService{
				TestService{Name: "nil-service"}.ReturningNilResult(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first and second services returning error but third returning success should return third result": {
			services: []TestService{
				TestService{Name: "error-service-1"}.Failing(),
				TestService{Name: "error-service-2"}.Failing(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first and second services returning nil but third returning success should return third result": {
			services: []TestService{
				TestService{Name: "nil-service-1"}.ReturningNilResult(),
				TestService{Name: "nil-service-2"}.ReturningNilResult(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first service returning error, second returning nil, third returning success should return third result": {
			services: []TestService{
				TestService{Name: "error-service"}.Failing(),
				TestService{Name: "nil-service"}.ReturningNilResult(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"first service returning nil, second returning error, third returning success should return third result": {
			services: []TestService{
				TestService{Name: "nil-service"}.ReturningNilResult(),
				TestService{Name: "error-service"}.Failing(),
				TestService{Name: "success-service"}.Successful(),
				TestService{Name: "should-not-be-called"}.ShouldNotBeCalled(),
			},
			expectedResult:   &Result{200, "success"},
			errorExpectation: assert.NoError,
		},
		"empty queue should return error": {
			services:         []TestService{},
			expectedResult:   nil,
			errorExpectation: assert.Error,
		},
	}

	for name, test := range tests {
		t.Run(name+": Queue", func(t *testing.T) {
			// given:
			services := slices.Map(test.services, func(s TestService) *servicequeue.Service[*Result] {
				service := s.NewTest(t)
				return servicequeue.NewService(service.Name, service.Do)
			})

			// and:
			queue := servicequeue.NewQueue(
				logging.NewTestLogger(t),
				"Do",
				services...,
			)

			// when:
			r, err := queue.OneByOne(context.Background())

			// then:
			test.errorExpectation(t, err)
			assert.Equal(t, test.expectedResult, r)
		})

		t.Run(name+": Queue1", func(t *testing.T) {
			// given:
			services := slices.Map(test.services, func(s TestService) *servicequeue.Service1[string, *Result] {
				service := s.NewTest(t)
				return servicequeue.NewService1(service.Name, service.Do1)
			})

			// and:
			queue := servicequeue.NewQueue1(
				logging.NewTestLogger(t),
				"Do1",
				services...,
			)

			// when:
			r, err := queue.OneByOne(context.Background(), secondArgument)

			// then:
			test.errorExpectation(t, err)
			assert.Equal(t, test.expectedResult, r)
		})

		t.Run(name+": Queue2", func(t *testing.T) {
			// given:
			services := slices.Map(test.services, func(s TestService) *servicequeue.Service2[string, int, *Result] {
				service := s.NewTest(t)
				return servicequeue.NewService2(service.Name, service.Do2)
			})

			// and:
			queue := servicequeue.NewQueue2(
				logging.NewTestLogger(t),
				"Do3",
				services...,
			)

			// when:
			r, err := queue.OneByOne(context.Background(), secondArgument, thirdArgument)

			// then:
			test.errorExpectation(t, err)
			assert.Equal(t, test.expectedResult, r)
		})

		t.Run(name+": Queue3", func(t *testing.T) {
			// given:
			services := slices.Map(test.services, func(s TestService) *servicequeue.Service3[string, int, bool, *Result] {
				service := s.NewTest(t)
				return servicequeue.NewService3(service.Name, service.Do3)
			})

			// and:
			queue := servicequeue.NewQueue3(
				logging.NewTestLogger(t),
				"Do3",
				services...,
			)

			// when:
			r, err := queue.OneByOne(context.Background(), secondArgument, thirdArgument, fourthArgument)

			// then:
			test.errorExpectation(t, err)
			assert.Equal(t, test.expectedResult, r)
		})
	}
}

type Result struct {
	StatusCode int
	Status     string
}

type TestService struct {
	Name         string
	t            testing.TB
	createResult func() (*Result, error)
}

func (s TestService) Successful() TestService {
	s.createResult = func() (*Result, error) {
		return &Result{200, "success"}, nil
	}
	return s
}

func (s TestService) Failing() TestService {
	s.createResult = func() (*Result, error) {
		return nil, errors.New("some error occurred")
	}
	return s
}

func (s TestService) ReturningNilResult() TestService {
	s.createResult = func() (*Result, error) {
		return nil, nil
	}
	return s
}

func (s TestService) Panicking() TestService {
	s.createResult = func() (*Result, error) {
		panic("some panic occurred")
	}
	return s
}

func (s TestService) ShouldNotBeCalled() TestService {
	s.createResult = func() (*Result, error) {
		s.t.Fatalf("service %s shouldn't be called, but was.", s.Name)
		return nil, nil
	}
	return s
}

func (s TestService) NewTest(t testing.TB) *TestService {
	s.t = t
	return &s
}

func (s *TestService) Do(ctx context.Context) (*Result, error) {
	assert.NotNil(s.t, ctx, "expect to receive non-nil context as 1st argument")
	return s.createResult()
}

func (s *TestService) Do1(ctx context.Context, str string) (*Result, error) {
	assert.Equal(s.t, secondArgument, str, "expect to receive %#v as 2nd argument", secondArgument)
	return s.Do(ctx)
}

func (s *TestService) Do2(ctx context.Context, str string, i int) (*Result, error) {
	assert.Equal(s.t, thirdArgument, i, "expect to receive %#v as 3rd argument", thirdArgument)
	return s.Do1(ctx, str)
}

func (s *TestService) Do3(ctx context.Context, str string, i int, boolean bool) (*Result, error) {
	assert.Equal(s.t, fourthArgument, boolean, "expect to receive %#v as 4th argument", fourthArgument)
	return s.Do2(ctx, str, i)
}
