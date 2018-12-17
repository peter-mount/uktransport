package publishmq

import (
  "errors"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/rabbitmq"
  "github.com/peter-mount/uktransport/lib"
  "os"
)

type PublishMQ struct {
  fileReader   *lib.FileReader
  rabbitMQ      rabbitmq.RabbitMQ

  url          *string
  exchange     *string
  routingKey   *string
}

func (a *PublishMQ) Name() string {
  return "publishmq"
}

func (a *PublishMQ) Init( k *kernel.Kernel ) error {
  a.url = flag.String( "u", "", "AMQP url to connect to broker" )
  a.exchange = flag.String( "exchange", "amq.topic", "Exchange to connect to" )
  a.routingKey = flag.String( "r", "", "Routing key to submit messages to" )

  service, err := k.AddService( &lib.FileReader{} )
  if err != nil {
    return err
  }
  a.fileReader = (service).(*lib.FileReader)

  a.fileReader.RecordProcessor = func(b []byte ) error {
    return a.rabbitMQ.Publish( *a.routingKey, b )
  }

  return nil
}

func (a *PublishMQ) PostInit() error {
  if *a.url == "" {
    *a.url = os.Getenv( "AMQP_URL" )
  }
  if *a.url == "" {
    return errors.New( "No amqp url via either -a or AMQP_URL" )
  }

  if len( flag.Args() ) == 0 {
    return errors.New( "No files supplied to parse" )
  }

  return nil
}

func (a *PublishMQ) Start() error {
  a.rabbitMQ.Url = *a.url
  a.rabbitMQ.Exchange = *a.exchange
  a.rabbitMQ.ConnectionName = "uktransport-publishmq"

  return a.rabbitMQ.Connect()
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
