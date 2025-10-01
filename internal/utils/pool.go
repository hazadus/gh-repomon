package utils

import (
	"context"
	"sync"
)

// ProcessInParallel processes items in parallel using a worker pool
// maxWorkers controls the maximum number of concurrent workers
// processor is called for each item and should return an error if processing fails
func ProcessInParallel[T any](items []T, maxWorkers int, processor func(T) error) error {
	if len(items) == 0 {
		return nil
	}

	// Limit workers to number of items if fewer items than workers
	if maxWorkers > len(items) {
		maxWorkers = len(items)
	}

	// Create channels for work distribution
	itemChan := make(chan T, len(items))
	errChan := make(chan error, len(items))

	// Send all items to the channel
	for _, item := range items {
		itemChan <- item
	}
	close(itemChan)

	// Start workers
	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wg.Done()
			for item := range itemChan {
				if err := processor(item); err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessInParallelWithContext processes items in parallel with context support
// The context can be used to cancel processing
func ProcessInParallelWithContext[T any](ctx context.Context, items []T, maxWorkers int, processor func(context.Context, T) error) error {
	if len(items) == 0 {
		return nil
	}

	if maxWorkers > len(items) {
		maxWorkers = len(items)
	}

	itemChan := make(chan T, len(items))
	errChan := make(chan error, 1)

	for _, item := range items {
		itemChan <- item
	}
	close(itemChan)

	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wg.Done()
			for item := range itemChan {
				select {
				case <-ctx.Done():
					return
				default:
					if err := processor(ctx, item); err != nil {
						select {
						case errChan <- err:
						default:
						}
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	select {
	case err := <-errChan:
		return err
	default:
		return ctx.Err()
	}
}
