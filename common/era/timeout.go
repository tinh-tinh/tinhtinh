package era

import (
	"context"
	"errors"
	"time"
)

func TimeoutFunc(duration time.Duration, task func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- task(ctx)
	}()

	select {
	case err := <-done:
		// Task completed before timeout
		return err
	case <-ctx.Done():
		// Timeout occurred
		return errors.New("function execution timed out")
	}
}
