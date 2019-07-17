package dbrest

import (
	"flag"
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/golib/kernel/cron"
	"github.com/peter-mount/golib/rest"
	"path/filepath"
)

type DBRest struct {
	configFile *string
	config     Config
	cron       *cron.CronService
	rest       *rest.Server
}

func (a *DBRest) Name() string {
	return "DBRest"
}

func (a *DBRest) Init(k *kernel.Kernel) error {
	a.configFile = flag.String("c", "", "The config file to use")

	service, err := k.AddService(&cron.CronService{})
	if err != nil {
		return err
	}
	a.cron = (service).(*cron.CronService)

	service, err = k.AddService(&rest.Server{})
	if err != nil {
		return err
	}
	a.rest = (service).(*rest.Server)

	return nil
}

func (a *DBRest) PostInit() error {
	if *a.configFile == "" {
		*a.configFile = "config.yaml"
	}

	filename, err := filepath.Abs(*a.configFile)
	if err != nil {
		return err
	}

	return a.config.unmarshal(nil, filename)
}

func (a *DBRest) Start() error {

	return a.config.start(a)
}
