package publishmq

import (
  "errors"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/uktransport/lib"
)

type PublishMQ struct {
  fileReader   *lib.FileReader
}

func (a *PublishMQ) Name() string {
  return "publishmq"
}

func (a *PublishMQ) Init( k *kernel.Kernel ) error {

  service, err := k.AddService( &lib.FileReader{} )
  if err != nil {
    return err
  }
  a.fileReader = (service).(*lib.FileReader)

  return nil
}

func (a *PublishMQ) PostInit() error {
  if len( flag.Args() ) == 0 {
    return errors.New( "No files supplied to parse" )
  }

  return nil
}

func (a *PublishMQ) Run() error {
  for _, s := range flag.Args() {
    err := a.fileReader.ReadFile( s )
    if err != nil {
      return err
    }
  }
  return nil
}
