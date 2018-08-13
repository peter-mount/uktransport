package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/uktransport/naptanimport"
  "log"
)

func main() {
  err := kernel.Launch( &naptanimport.NaptanImport{} )
  if err != nil {
    log.Fatal( err )
  }
}
