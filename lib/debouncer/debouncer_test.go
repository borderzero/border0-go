package debouncer

import (
	"sync/atomic"
	"testing"
	"time"
)

// helper to receive exactly one value from ch within d, failing the test otherwise.
func recvOne[T any](t *testing.T, ch <-chan T, within time.Duration) T {
	t.Helper()

	select {
	case v := <-ch:
		return v
	case <-time.After(within):
		t.Fatalf("timed out waiting for value within %v", within)
		var zero T
		return zero
	}
}

// helper to assert no value received from ch within d.
func recvNone[T any](t *testing.T, ch <-chan T, within time.Duration) {
	t.Helper()

	select {
	case v := <-ch:
		t.Fatalf("expected no value, but got: %v", v)
	case <-time.After(within):
		// good
	}
}

func TestDebounce_LastCallWins(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	// Use small debounce window, no max wait.
	d := New(
		func(v int) { calls <- v },
		WithDebounceTime(50*time.Millisecond),
	)
	defer d.Close()

	// Rapid calls; only the LAST should be invoked after debounce elapses.
	d.Do(1)
	time.Sleep(10 * time.Millisecond)
	d.Do(2)
	time.Sleep(10 * time.Millisecond)
	d.Do(3)

	// Before debounce elapses, nothing should fire.
	recvNone(t, calls, 30*time.Millisecond)

	// After debounce, exactly one call with last value.
	got := recvOne(t, calls, 150*time.Millisecond)
	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}

	// And no extra calls afterwards.
	recvNone(t, calls, 100*time.Millisecond)
}

func TestMaxWait_TriggersFlushWithLatestArg(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	// Make debounce very long so it wouldn't fire; rely on maxWait instead.
	d := New(
		func(v int) { calls <- v },
		WithDebounceTime(1*time.Second),
		WithMaxWaitTime(120*time.Millisecond),
	)
	defer d.Close()

	d.Do(10)
	time.Sleep(60 * time.Millisecond)
	d.Do(20) // arrives before maxWait; latest arg should be 20

	// Expect maxWait to trigger a flush ~120ms after the first Do.
	got := recvOne(t, calls, 400*time.Millisecond)
	if got != 20 {
		t.Fatalf("expected 20 due to latest-arg policy at maxWait, got %d", got)
	}

	// No extra call should occur since timers are reset/cleared upon Flush.
	recvNone(t, calls, 200*time.Millisecond)
}

func TestChannelBackpressure_LatestWinsUnderLoad(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	d := New(
		func(v int) { calls <- v },
		WithDebounceTime(60*time.Millisecond),
	)
	defer d.Close()

	// Flood with many calls quickly; the debouncer's single-slot channel drops older enqueued messages.
	for i := range 100 {
		d.Do(i)
		// very small spacing to try to keep channel pressured
		if i%5 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}

	// After the debounce window, we should get exactly one call — the latest value.
	got := recvOne(t, calls, 300*time.Millisecond)
	if got != 99 {
		t.Fatalf("expected last value 99, got %d", got)
	}

	// No additional calls after the first flush.
	recvNone(t, calls, 150*time.Millisecond)
}

func TestFlush_Immediate(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	d := New(
		func(v int) { calls <- v },
		WithDebounceTime(200*time.Millisecond),
	)
	defer d.Close()

	d.Do(42)

	// Flush right away; should call immediately (not waiting for debounce).
	d.Flush()

	got := recvOne(t, calls, 10*time.Millisecond)
	if got != 42 {
		t.Fatalf("expected 42 after immediate Flush, got %d", got)
	}

	// Ensure no residual timers cause a duplicate call.
	recvNone(t, calls, 300*time.Millisecond)
}

func TestFlush_TwiceAfterDo_SecondIsNoOp(t *testing.T) {
	calls := make(chan int, 2)
	defer close(calls)

	d := New(
		func(v int) { calls <- v },
		WithDebounceTime(250*time.Millisecond),
	)
	defer d.Close()

	d.Do(99)

	// First flush should invoke once.
	d.Flush()
	got := recvOne(t, calls, 150*time.Millisecond)
	if got != 99 {
		t.Fatalf("expected 99 on first flush, got %d", got)
	}

	// Second flush immediately after should be a no-op.
	d.Flush()
	recvNone(t, calls, 150*time.Millisecond)
}

func TestClose_FlushesPendingAndIsIdempotent(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	var closedCalls atomic.Int32
	d := New(
		func(v int) {
			// track invocations just in case; also send to channel for assertion
			closedCalls.Add(1)
			calls <- v
		},
		WithDebounceTime(250*time.Millisecond),
	)

	d.Do(7)

	// Close should drain channel processor, wait, then Flush pending call.
	d.Close()

	got := recvOne(t, calls, 200*time.Millisecond)
	if got != 7 {
		t.Fatalf("expected pending value 7 to be flushed on Close, got %d", got)
	}

	// Closing again should not panic nor invoke again.
	d.Close()
	if closedCalls.Load() != 1 {
		t.Fatalf("expected exactly 1 invocation, got %d", closedCalls.Load())
	}
}

func TestDoAfterClose_IsIgnored(t *testing.T) {
	calls := make(chan int, 10)
	defer close(calls)

	d := New(func(v int) { calls <- v }, WithDebounceTime(50*time.Millisecond))
	d.Close()

	// These should be ignored silently.
	d.Do(1)
	d.Do(2)

	// Flush after close should do nothing as there is no pending work.
	d.Flush()

	recvNone(t, calls, 150*time.Millisecond)
}
