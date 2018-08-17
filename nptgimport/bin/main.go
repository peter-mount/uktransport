package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/uktransport/nptgimport"
  "log"
)

func main() {
  err := kernel.Launch( &kernel.MemUsage{}, &nptgimport.NptgImport{} )
  if err != nil {
    log.Fatal( err )
  }
}
