package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/uktransport/publishmq"
  "log"
)

func main() {
  err := kernel.Launch( &kernel.MemUsage{}, &publishmq.PublishMQ{} )
  if err != nil {
    log.Fatal( err )
  }
}
