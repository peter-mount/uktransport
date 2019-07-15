package dbrest

import (
	"flag"
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/golib/kernel/cron"
	"github.com/peter-mount/golib/kernel/db"
	"github.com/peter-mount/golib/rest"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type DBRest struct {
	configFile *string
	config     Config
	cron       *cron.CronService
	db         *db.DBService
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

	service, err = k.AddService(&db.DBService{})
	if err != nil {
		return err
	}
	a.db = (service).(*db.DBService)

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

	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(in, &a.config)
	if err != nil {
		return err
	}

	a.db.SetDB(a.config.DB.PostgresUri).
		MaxOpen(a.config.DB.MaxOpen).
		MaxIdle(a.config.DB.MaxIdle).
		MaxLifetime(time.Second * time.Duration(a.config.DB.MaxLifetime))

	return nil
}

func (a *DBRest) Start() error {

	for _, handler := range a.config.Rest {
		handler.init(a.db, a.rest)
		log.Println(handler.Path, handler.sql)
	}

	for i, handler := range a.config.Rest {
		log.Println(i, handler.Path, handler.sql)
	}

	return nil
}
