// gfurture is a simple implementation of a future/promise pattern in Go.
//
// It allows you to create a future value that can be resolved later.
// This is useful for asynchronous programming, where you may want to perform
// some computation in the background and get the result later.
//
// The package provides a Future type that can be used to create and manage
// future values. You can create a Future using the NewFuture function, and
// resolve it using the Resolve method.
// The Await method blocks until the future is resolved, and returns the value
// and error.
// The Then method allows you to execute a consumer function with the value
// and error of the future once it is resolved. This can be useful for chaining
// multiple futures together, or for performing some action once the future is resolved.
//
// The package is designed to be simple and easy to use, with a minimal API.
// It is not intended to be a full-featured implementation of the future/promise
// pattern, but rather a lightweight alternative for Go developers who want
// to use this pattern in their code. The package is also designed to be
// thread-safe, so you can use it in concurrent programs without worrying.
package gfuture

import (
	"context"
)

type payload[T any] struct {
	val T     // The value of the payload.
	err error // The error associated with the payload, if any.
}

// Future is a generic type representing a future value that will be available later.
type Future[T any] chan payload[T]

// NewFuture creates and returns a new Future instance.
func NewFuture[T any]() Future[T] {
	return make(chan payload[T])
}

// Async creates a Future and executes the provided function asynchronously.
// The result of the function is resolved into the Future.
func Async[T any](provider func() (T, error)) Future[T] {
	f := NewFuture[T]()
	go func() {
		f.Resolve(provider())
	}()
	return f
}

func (f Future[T]) sendAndClose(p payload[T]) {
	f <- p
	close(f)
}

// Resolve sets the value and error of the Future and closes it.
func (f Future[T]) Resolve(value T, err error) {
	f.sendAndClose(payload[T]{val: value, err: err})
}

// Value sets the value of the Future and closes it.
func (f Future[T]) Value(value T) {
	f.sendAndClose(payload[T]{val: value})
}

// Error sets the error of the Future and closes it.
func (f Future[T]) Error(err error) {
	f.sendAndClose(payload[T]{err: err})
}

// Await waits for the Future to resolve and returns the value and error.
func (f Future[T]) Await(ctx context.Context) (T, error) {
	select {
	case payload := <-f:
		return payload.val, payload.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Then executes the provided consumer function with the value and error of the Future once resolved.
func (f Future[T]) Then(ctx context.Context, consumer func(T, error)) {
	go func() {
		consumer(f.Await(ctx))
	}()
}
