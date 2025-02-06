package era_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common/era"
)

func Test_Timeout(t *testing.T) {
	longTask := func(ctx context.Context) error {
		fmt.Println("Long task started")
		select {
		case <-time.After(3 * time.Second): // Simulate work that takes 3 seconds
			fmt.Println("Long task completed")
			return nil
		case <-ctx.Done():
			fmt.Println("Long task was canceled")
			return ctx.Err()
		}
	}

	shortTask := func(ctx context.Context) error {
		fmt.Println("Short task started")
		time.Sleep(1 * time.Second) // Simulate a short task
		fmt.Println("Short task completed")
		return nil
	}

	parameterizedTask := func(param string, delay time.Duration) func(context.Context) error {
		return func(ctx context.Context) error {
			fmt.Printf("Task with param '%s' started\n", param)
			select {
			case <-time.After(delay):
				fmt.Printf("Task with param '%s' completed\n", param)
				return nil
			case <-ctx.Done():
				fmt.Printf("Task with param '%s' was canceled\n", param)
				return ctx.Err()
			}
		}
	}

	fmt.Println("Running long-running task with 2-second timeout:")
	err := era.TimeoutFunc(2*time.Second, longTask)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\nRunning short task with 2-second timeout:")
	err = era.TimeoutFunc(2*time.Second, shortTask)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\nRunning parameterized task with 2-second timeout:")
	task := parameterizedTask("exampleParam", 3*time.Second)
	err = era.TimeoutFunc(2*time.Second, task)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
