package golang

import (
  "playground/internal/log"
  "sync"
)

const (
  workers = 10
  tasks   = 5
)

var (
  jobs    = make(chan int)
  results = make(chan int)
)

// 1. Simultaneously spawn workers (W), a producer (P) & a consumer (C)
// 2. W waits for jobs <-|
// 3.      P adds jobs ->|
// 4. W adds processed jobs to results ->|
// 5.              C waits for results <-|

func Channel() {
  var wg sync.WaitGroup

  for i := 0; i < workers; i++ {
    wg.Add(1)
    go fanOut(&wg, i)
  }

  go producer()
  go consumer()
  fanIn(&wg)
}

func fanOut(wg *sync.WaitGroup, idx int) {
  defer wg.Done()

  // NOTE: If tasks < worker, non-chosen workers (runtime pseudo-randomly
  // selects a routine) will return immediately.
  for j := range jobs {
    log.Info("%d: %d", idx, j)
    results <- j << 1
  }
}

func producer() {
  for i := 0; i < tasks; i++ {
    jobs <- i
  }

  close(jobs)
}

func consumer() {
  for r := range results {
    log.Debug("read: %d", r)
  }
}

func fanIn(wg *sync.WaitGroup) {
  wg.Wait()
  log.Info("fan-out/fan-in finished")
  close(results)
}
