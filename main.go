package main

import (
  "os"
  "playground/internal/log"
  "playground/pkg/foo"
)

var modules map[string]func()

func init() {
  modules = make(map[string]func())
  modules["foo"] = foo.Run

  log.Debug("%d module(s) initialised", len(modules))
}

func main() {
  args := os.Args[1:]
  if len(args) != 1 {
    log.Error("too many/few arguments")
    os.Exit(1)
  }

  module := args[0]
  fn, found := modules[module]
  if !found {
    log.Error("unknown module: %s", module)
    os.Exit(1)
  }

  log.Info("-> %s", module)
  fn()
}
