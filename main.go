package main

import (
  "fmt"
  "os"
  "playground/internal/log"
  "playground/pkg/golang"
  "playground/pkg/structure"
)

var modules map[string]func()

func init() {
  modules = make(map[string]func())
  modules["atomic"] = golang.Atomic
  modules["channel"] = golang.Channel
  modules["context"] = golang.Context
  modules["linkedlist"] = structure.LinkedList

  log.Debug("%d module(s) ready", len(modules))
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
