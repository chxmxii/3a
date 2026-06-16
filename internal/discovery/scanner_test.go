package discovery

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryWithBackoff_SuccessOnFirstTry(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	callCount := 0
	err := retryWithBackoff(context.Background(), cfg, func() error {
		callCount++
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected fn to be called once, got %d calls", callCount)
	}
}

func TestRetryWithBackoff_SuccessAfterFailures(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	// Succeed on the 3rd call (after 2 failures).
	callCount := 0
	err := retryWithBackoff(context.Background(), cfg, func() error {
		callCount++
		if callCount < 3 {
			return errors.New("transient error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if callCount != 3 {
		t.Fatalf("expected fn to be called 3 times, got %d", callCount)
	}
}

func TestRetryWithBackoff_SuccessAfterOneFailure(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	// Succeed on the 2nd call (after 1 failure).
	callCount := 0
	err := retryWithBackoff(context.Background(), cfg, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("transient error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected fn to be called 2 times, got %d", callCount)
	}
}

func TestRetryWithBackoff_AllRetriesExhausted(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	callCount := 0
	expectedErr := errors.New("persistent error")
	err := retryWithBackoff(context.Background(), cfg, func() error {
		callCount++
		return expectedErr
	})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected error %q, got %q", expectedErr, err)
	}
	// 1 initial call + 3 retries = 4 total calls
	if callCount != 4 {
		t.Fatalf("expected fn to be called 4 times (1 initial + 3 retries), got %d", callCount)
	}
}

func TestRetryWithBackoff_ContextCancelledDuringBackoff(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond, // Long enough that we can cancel during it
		BackoffFactor:  2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())

	var callCount int32
	go func() {
		// Wait for the first call to complete, then cancel.
		for atomic.LoadInt32(&callCount) < 1 {
			time.Sleep(5 * time.Millisecond)
		}
		cancel()
	}()

	err := retryWithBackoff(ctx, cfg, func() error {
		atomic.AddInt32(&callCount, 1)
		return errors.New("will fail")
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

func TestRetryWithBackoff_ExponentialBackoffPattern(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 50 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	// Track timestamps of each call to verify exponential backoff.
	var timestamps []time.Time
	_ = retryWithBackoff(context.Background(), cfg, func() error {
		timestamps = append(timestamps, time.Now())
		return errors.New("keep failing")
	})

	if len(timestamps) != 4 {
		t.Fatalf("expected 4 timestamps, got %d", len(timestamps))
	}

	// Expected backoff intervals: 50ms, 100ms, 200ms
	// Allow 30ms tolerance for scheduling jitter.
	expectedBackoffs := []time.Duration{
		50 * time.Millisecond,  // initial * 2^0
		100 * time.Millisecond, // initial * 2^1
		200 * time.Millisecond, // initial * 2^2
	}
	tolerance := 30 * time.Millisecond

	for i := 0; i < 3; i++ {
		actual := timestamps[i+1].Sub(timestamps[i])
		expected := expectedBackoffs[i]
		if actual < expected-tolerance || actual > expected+tolerance {
			t.Errorf("backoff %d: expected ~%v, got %v", i, expected, actual)
		}
	}
}

func TestRetryWithBackoff_ZeroRetries(t *testing.T) {
	cfg := RetryConfig{
		MaxRetries:     0,
		InitialBackoff: 10 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	callCount := 0
	err := retryWithBackoff(context.Background(), cfg, func() error {
		callCount++
		return errors.New("fails immediately")
	})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if callCount != 1 {
		t.Fatalf("expected fn to be called once with MaxRetries=0, got %d", callCount)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()

	if cfg.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", cfg.MaxRetries)
	}
	if cfg.InitialBackoff != 1*time.Second {
		t.Errorf("expected InitialBackoff=1s, got %v", cfg.InitialBackoff)
	}
	if cfg.BackoffFactor != 2.0 {
		t.Errorf("expected BackoffFactor=2.0, got %f", cfg.BackoffFactor)
	}
}
