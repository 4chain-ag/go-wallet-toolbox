package servicequeue

import "fmt"

// ServiceExecutor is a struct that holds a name and Service that implements interface T
type ServiceExecutor[T any] struct {
	Name    string
	Service T
}

// ServicesQueue represents a queue of ServiceExecutors and its current index
type ServicesQueue[T any] struct {
	services []ServiceExecutor[T]
	index    int
}

// New will create a new ServicesQueue that holds services to execute method T
func New[T any](services ...ServiceExecutor[T]) *ServicesQueue[T] {
	return &ServicesQueue[T]{
		services: services,
		index:    0,
	}
}

// Current will return the current service in queue if there is any
func (sq *ServicesQueue[T]) Current() (ServiceExecutor[T], error) {
	if len(sq.services) == 0 {
		return ServiceExecutor[T]{}, fmt.Errorf("no services available")
	}
	return sq.services[sq.index], nil
}

func (sq *ServicesQueue[T]) Add(service ServiceExecutor[T]) *ServicesQueue[T] {
	sq.services = append(sq.services, service)
	return sq
}

// Next moves to the next service in the queue
func (sq *ServicesQueue[T]) Next() {
	if len(sq.services) > 0 {
		sq.index = (sq.index + 1) % len(sq.services)
	}
}

// Remove removes a service from the collection by its name.
func (sq *ServicesQueue[T]) Remove(name string) {
	var newServices []ServiceExecutor[T]
	for _, s := range sq.services {
		if s.Name != name {
			newServices = append(newServices, s)
		}
	}

	sq.services = newServices

	if sq.index >= len(sq.services) {
		sq.index = 0
	}
}

// Count returns number of services that the queue is holding at the moment
func (sq *ServicesQueue[T]) Count() int {
	return len(sq.services)
}
