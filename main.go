package main

import (
  "fmt"
  "os"
  "playground/internal/log"
  "playground/pkg/channel"
  "playground/pkg/context"
)

var modules map[string]func()

func init() {
  modules = make(map[string]func())
  modules["context"] = context.Run
  modules["channel"] = channel.Run

  log.Debug("%d module(s) initialised", len(modules))
}

func main() {
  args := os.Args[1:]
  if len(args) != 1 {
    log.Abort("too many/few arguments")
  }

  module := args[0]
  fn, found := modules[module]
  if !found {
    log.Abort("unknown module: %s", module)
  }

  fmt.Println()
  fn()
}
