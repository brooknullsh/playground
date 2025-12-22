package log

import (
  "fmt"
  "os"
  "strings"
)

const (
  RESET = "\033[0m"
  DEBUG = "\033[0;34m"
  INFO  = "\033[0;32m"
  WARN  = "\033[0;33m"
  ERROR = "\033[0;31m"
)

func Debug(format string, args ...any) {
  fmt.Fprintf(os.Stdout, coloured(DEBUG, format), args...)
}

func Info(format string, args ...any) {
  fmt.Fprintf(os.Stdout, coloured(INFO, format), args...)
}

func Warn(format string, args ...any) {
  fmt.Fprintf(os.Stderr, coloured(WARN, format), args...)
}

func Error(format string, args ...any) {
  fmt.Fprintf(os.Stderr, coloured(ERROR, format), args...)
}

func coloured(colour, format string) string {
  var key string
  switch colour {
  case DEBUG:
    key = "[D] "
  case INFO:
    key = "[I] "
  case WARN:
    key = "[W] "
  case ERROR:
    key = "[E] "
  default:
    key = "[D] "
  }

  var builder strings.Builder
  builder.WriteString(colour)
  builder.WriteString(key)
  builder.WriteString(RESET)
  builder.WriteString(format)
  builder.WriteString("\n")

  return builder.String()
}
