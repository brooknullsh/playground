package atomic

import (
  "playground/internal/log"
  "sync"
  "sync/atomic"
)

// Each goroutine increments a shared counter. Since the actual work is very
// light, the contention (data mutated, cache line invalidated) is high.
const (
  steps = 10
  iters = 1000
)

func Run() {
  var (
    mu      sync.Mutex
    muCount int
    wg      sync.WaitGroup
  )

  // NOTE: A mutex is better here as the cache line is only invalidated on lock
  // and unlock while all other routines block (may even be descheduled) and
  // stop contending.
  for i := 0; i < steps; i++ {
    wg.Add(1)
    go mutexInc(&wg, &mu, &muCount)
  }

  var atCount atomic.Int32

  // NOTE: All cores repeatedly try and access the cache line while it's being
  // invalidated, leading to coherence misses.
  for i := 0; i < steps; i++ {
    wg.Add(1)
    go atomicInc(&wg, &atCount)
  }

  wg.Wait()
  log.Info("mutex: %d", muCount)
  log.Info("atomic: %d", atCount.Load())
}

func mutexInc(wg *sync.WaitGroup, mu *sync.Mutex, count *int) {
  defer wg.Done()

  // Each iteration gains ownership of the cache line. All other routines block,
  // the next iteration, another routine may have ownership.
  for i := 0; i < iters; i++ {
    mu.Lock()
    *count++
    mu.Unlock()
  }
}

func atomicInc(wg *sync.WaitGroup, count *atomic.Int32) {
  defer wg.Done()

  // No blocking for each iteration, performs an atomic RMW (Read-Modify-Write)
  // which forces invalidation and other accesses must retry.
  for i := 0; i < iters; i++ {
    count.Add(1)
  }
}
