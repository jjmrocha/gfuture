# gfuture

`gfuture` is a lightweight implementation of the future/promise pattern in Go. It allows you to create and manage future values that can be resolved later, making it easier to work with asynchronous programming in Go.

## Features

- Create future values using `NewFuture` or `Async`.
- Resolve futures with values or errors.
- Await the resolution of a future with optional context cancellation.
- Chain actions using the `Then` method.

## Installation

To use `gfuture`, simply import it into your Go project:

```go
import "path/to/gfuture"
```

## Usage

### Creating and Resolving a Future

```go
package main

import (
	"context"
	"fmt"
	"time"
	"gfuture"
)

func main() {
	// Create a new Future
	future := gfuture.NewFuture[int]()

	// Resolve the Future asynchronously
	go func() {
		time.Sleep(1 * time.Second)
		future.Resolve(42, nil)
	}()

	// Await the result
	ctx := context.Background()
	value, err := future.Await(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value) // Output: Value: 42
	}
}
```

### Using `Async` for Simpler Asynchronous Execution

```go
package main

import (
	"context"
	"fmt"
	"gfuture"
)

func main() {
	// Create a Future using Async
	future := gfuture.Async(func() (int, error) {
		return 42, nil
	})

	// Await the result
	ctx := context.Background()
	value, err := future.Await(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value) // Output: Value: 42
	}
}
```

### Chaining Actions with `Then`

```go
package main

import (
	"context"
	"fmt"
	"gfuture"
)

func main() {
  ctx := context.Background()
  
	// Create a Future using Async
	future := gfuture.Async(func() (int, error) {
		return 42, nil
	}).Then(ctx, func(value int, err error) {
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Value:", value) // Output: Value: 42
		}
	})

	// Wait to ensure the program doesn't exit prematurely
	select {}
}
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
