package servicequeue

import "context"

// Service represents a named service function accepting only context.Context as an argument.
type Service[R any] struct {
	name    string
	service func(context.Context) (R, error)
}

// NewService creates a new Service instance with given name.
func NewService[R any](name string, service func(context.Context) (R, error)) *Service[R] {
	return &Service[R]{
		name:    name,
		service: service,
	}
}

// Name returns the name of the service.
// Used for logging and debugging purposes.
func (s *Service[R]) Name() string {
	return s.name
}

// Service1 represents a named service function accepting context.Context and one additional argument.
type Service1[A, R any] struct {
	name    string
	service func(context.Context, A) (R, error)
}

// NewService1 creates a new Service1 instance with given name.
func NewService1[A, R any](name string, service func(context.Context, A) (R, error)) *Service1[A, R] {
	return &Service1[A, R]{
		name:    name,
		service: service,
	}
}

// Name returns the name of the service.
// Used for logging and debugging purposes.
func (s *Service1[A, R]) Name() string {
	return s.name
}

// Service2 represents a named service function accepting context.Context and two additional arguments.
// It is used for services that require two parameters to perform their operations.
type Service2[A, B, R any] struct {
	name    string
	service func(context.Context, A, B) (R, error)
}

// NewService2 creates a new Service2 instance with given name.
func NewService2[A, B, R any](name string, service func(context.Context, A, B) (R, error)) *Service2[A, B, R] {
	return &Service2[A, B, R]{
		name:    name,
		service: service,
	}
}

// Name returns the name of the service.
// Used for logging and debugging purposes.
func (s *Service2[A, B, R]) Name() string {
	return s.name
}

// Service3 represents a named service function accepting context.Context and three additional arguments.
type Service3[A, B, C, R any] struct {
	name    string
	service func(context.Context, A, B, C) (R, error)
}

// NewService3 creates a new Service3 instance with given name.
func NewService3[A, B, C, R any](name string, service func(context.Context, A, B, C) (R, error)) *Service3[A, B, C, R] {
	return &Service3[A, B, C, R]{
		name:    name,
		service: service,
	}
}

// Name returns the name of the service.
// Used for logging and debugging purposes.
func (s *Service3[A, B, C, R]) Name() string {
	return s.name
}
