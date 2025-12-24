package golang

import (
  "context"
  "playground/internal/log"
  "sync"
  "time"
)

const (
  max  = 3 * time.Second
  work = 1 * time.Second
)

func Context() {
  ctx, cancel := context.WithTimeout(context.Background(), max)
  log.Debug("max: %v", max)

  var wg sync.WaitGroup
  wg.Add(1)

  go bar(ctx, &wg)
  if err := foo(ctx); err != nil {
    log.Abort("foo() failed: %v", err)
  }

  // NOTE: This would usually be deferred above, but it's here for the success
  // state of foo() where bar() is waiting.
  cancel()
  wg.Wait()
}

func foo(ctx context.Context) error {
  select {
  // NOTE: Done() is a receive-only channel typed by struct{} (zero bytes) that
  // internally is closed once the context is cancelled.
  case <-ctx.Done():
    return ctx.Err()
  case <-time.After(work):
    log.Debug("worked: %v", work)
    return nil
  }
}

func bar(ctx context.Context, wg *sync.WaitGroup) {
  defer wg.Done()

  // NOTE: This will run for both success and failure cases as the context is
  // cancelled either way. The abort for the failure isn't quick enough to
  // prevent this from running.
  <-ctx.Done()
  log.Debug("routine finished")
}
