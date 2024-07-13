package main

import (
	"fmt"
	"sync"
)

func main() {
  ch := make(chan int)

  wg := sync.WaitGroup{}
  wg.Add(2)

  go func() {
    defer wg.Done()
    for x := range ch {
      fmt.Printf("from chan 1: %d\n", x)
    }
  }()

  go func() {
    defer wg.Done()
    for x := range ch {
      fmt.Printf("from chan 2: %d\n", x)
    }
  }()

  for i := 0; i < 10; i++ {
    ch <- i
  }
  close(ch)

  wg.Wait()
}
