package discovery

import (
	"context"
	"time"
)

// retryWithBackoff retries fn up to cfg.MaxRetries times with exponential backoff.
// The backoff duration follows the pattern: initial * factor^attempt.
// Returns nil immediately if fn succeeds, returns ctx.Err() if context is
// cancelled during backoff, or returns the last error if all retries are exhausted.
func retryWithBackoff(ctx context.Context, cfg RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// If this was the last allowed attempt, don't wait—just return the error.
		if attempt == cfg.MaxRetries {
			break
		}

		// Calculate backoff: initial * factor^attempt
		backoff := time.Duration(float64(cfg.InitialBackoff) * pow(cfg.BackoffFactor, attempt))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			// Continue to next retry attempt.
		}
	}

	return lastErr
}

// pow computes base^exp for float64 base and integer exponent.
func pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
