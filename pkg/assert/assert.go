package assert

import (
	"log"
	"log/slog"
)

var assertData map[string]any = map[string]any{}

func Assert(truth bool, msg string) {
  if !truth {
    runAssert(msg)
  }
}

func NoError(err error, msg string) {
  if err != nil {
    slog.Error("no error encountered", err)
    runAssert(msg)
  }
}

// utils
func runAssert(msg string) {
  for k, v := range assertData {
    slog.Error("context", "key", k, "value", v)
  }
  log.Fatal(msg)
}

func AddAssertData(key string, value any) {
  assertData[key] = value
}

func RemoveAssertData(key string) {
  delete(assertData, key)
}
