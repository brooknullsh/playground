package context

import (
  "context"
  "playground/internal/log"
  "time"
)

const (
  wait = 2 * time.Second
  busy = 1 * time.Second
)

func Run() {
  // Create a new context with an empty parent, start a timer that cancels all
  // children using said context once expired.
  ctx, cancel := context.WithTimeout(context.Background(), wait)
  defer cancel() // Must always call cancel() to release resources.

  log.Debug("limit: %v", wait)
  log.Info("foo()<-")

  // cancel()
  // log.Debug("cancelling context before foo()")
  if err := foo(ctx); err != nil {
    // The cancellation reason can be checked for specifics:
    // ... errors.Is(err, context.Canceled)
    // ... errors.Is(err, context.DeadlineExceeded)

    // NOTE: Technically, Abort() calls os.Exit() which doesn't run deferred
    // functions so cancel() won't be called.
    log.Abort("foo() failed: %v", err)
  }

  log.Info("<-foo()")
}

func foo(ctx context.Context) error {
  select {
  // Mock some workload. If the passed context hasn't been cancelled, will block
  // until the time has passed.
  //
  // NOTE: If this blocked for longer than the passed context timeout, the
  // context will cancel.
  case <-time.After(busy):
    log.Debug("worked: %v", busy)
    return nil
  // Cancelled context, could have been explicitly cancelled or timed out.
  //
  // NOTE: Done() is just a receive-only channel typed by struct{} (zero bytes)
  // that internally is closed once the context is cancelled. Channel closure
  // satisfies this case just as well as a value received.
  case <-ctx.Done():
    return ctx.Err() // Error is non-nil until context is cancelled.
  }
}
