package gfuture

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewFuture(t *testing.T) {
	future := NewFuture[int]()
	if future == nil {
		t.Fatal("Expected a non-nil Future")
	}
}

func TestResolve(t *testing.T) {
	// given
	ctx := context.Background()
	future := NewFuture[int]()
	// when
	go future.Resolve(42, nil)
	value, err := future.Await(ctx)
	// then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != 42 {
		t.Fatalf("Expected value 42, got %v", value)
	}
}

func TestValue(t *testing.T) {
	// given
	ctx := context.Background()
	future := NewFuture[int]()
	// when
	go future.Value(42)
	value, err := future.Await(ctx)
	// then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != 42 {
		t.Fatalf("Expected value 42, got %v", value)
	}
}

func TestError(t *testing.T) {
	// given
	ctx := context.Background()
	expectedErr := errors.New("test error")
	future := NewFuture[int]()
	// when
	go future.Error(expectedErr)
	value, err := future.Await(ctx)
	// then
	if err != expectedErr {
		t.Fatalf("Expected error %v, got %v", expectedErr, err)
	}

	if value != 0 {
		t.Fatalf("Expected value 0, got %v", value)
	}
}

func TestAsync(t *testing.T) {
	// given
	ctx := context.Background()
	// when
	value, err := Async(func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}).Await(ctx)
	// then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != 42 {
		t.Fatalf("Expected value 42, got %v", value)
	}
}

func TestAwaitWithTimeout(t *testing.T) {
	// given
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	// when
	value, err := Async(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	}).Await(ctx)
	// then
	if err == nil {
		t.Fatal("Expected context deadline exceeded error, got nil")
	}

	if value != 0 {
		t.Fatalf("Expected value 0, got %v", value)
	}
}

func TestThen(t *testing.T) {
	// given
	ctx := context.Background()
	future := NewFuture[int]()
	// when
	Async(func() (int, error) {
		return 42, nil
	}).Then(ctx, func(value int, err error) {
		future.Resolve(value, err)
	})
	value, err := future.Await(ctx)
	// then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value != 42 {
		t.Fatalf("Expected value 42, got %v", value)
	}
}

func TestThenWithError(t *testing.T) {
	// given
	ctx := context.Background()
	expectedErr := errors.New("test error")
	future := NewFuture[int]()
	// when
	Async(func() (int, error) {
		return 0, expectedErr
	}).Then(ctx, func(value int, err error) {
		future.Resolve(value, err)
	})
	value, err := future.Await(ctx)
	// then
	if err != expectedErr {
		t.Fatalf("Expected error %v, got %v", expectedErr, err)
	}

	if value != 0 {
		t.Fatalf("Expected value 0, got %v", value)
	}
}
