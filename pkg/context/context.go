package context

import (
  "context"
  "playground/internal/log"
  "sync"
  "time"
)

const (
  wait  = 3 * time.Second
  busy  = 3 * time.Second
  pause = 500 * time.Millisecond
)

func Run() {
  // Create a new context with an empty parent, start a timer that cancels all
  // children using said context once expired.
  ctx, cancel := context.WithTimeout(context.Background(), wait)

  log.Debug("limit: %v", wait)
  var wg sync.WaitGroup
  wg.Add(1)

  // Pass the context to a goroutine. The goroutine will work alongside the
  // below synchronous call and check for the context's cancellation.
  go bar(ctx, &wg)

  if err := foo(ctx); err != nil {
    // The cancellation reason can be checked for specifics:
    // ... errors.Is(err, context.DeadlineExceeded)

    // NOTE: Abort() calls os.Exit() which doesn't run deferred functions so
    // cancel() wouldn't be called.
    log.Abort("foo() failed: %v", err)
  }

  // Cancel the context at this point, instead of waiting for the timeout. The
  // goroutine should then resolve.
  //
  // NOTE: This would usually be deferred above.
  cancel()
  wg.Wait()
}

func foo(ctx context.Context) error {
  log.Info("foo()<-")

  select {
  // Cancelled context, could have been explicitly cancelled or timed out.
  //
  // NOTE: Done() is just a receive-only channel typed by struct{} (zero bytes)
  // that internally is closed once the context is cancelled. Channel closure
  // satisfies this case just as well as a value received.
  case <-ctx.Done():
    return ctx.Err() // Error is non-nil until context is cancelled.
  // Mock some workload. If the passed context hasn't been cancelled, will block
  // until the time has passed.
  //
  // NOTE: If this blocked for longer than the passed context timeout, the
  // context will cancel.
  case <-time.After(busy):
    log.Debug("worked: %v", busy)
    log.Info("<-foo()")
    return nil
  }
}

func bar(ctx context.Context, wg *sync.WaitGroup) {
  defer wg.Done()
  log.Info("bar()<-")

  // Block until the context is cancelled. The sync.WaitGroup allows the "main
  // thread" to wait for the cancellation case to match.
  //
  // NOTE: If the "busy" and "wait" times are equal, this may run as the "busy"
  // select case runs.
  <-ctx.Done()
  log.Info("<-bar()")
}
