package dbrest

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Config struct {
	prefix   string
	DB       *DB               `yaml:"db"`
	Cronjobs []CronJob         `yaml:"cron"`
	Rest     []*RestHandler    `yaml:"rest"`
	Imports  map[string]string `yaml:"import"`
	imports  []*Config
	handlers []*RestHandler
}

func (c *Config) unmarshal(parent *Config, filename string) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	log.Println("Loading", filename)

	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(in, c)
	if err != nil {
		return err
	}

	if c.DB == nil {
		if parent == nil {
			return errors.New("Database is mandatory for the root config.yaml")
		}
		c.DB = parent.DB
	}

	if len(c.Imports) > 0 {

		base := filepath.Dir(filename)

		for prefix, path := range c.Imports {
			importFileName := filepath.Join(base, path)
			importFileName, err := filepath.Abs(importFileName)
			if err != nil {
				return err
			}

			ci := &Config{prefix: prefix}
			err = ci.unmarshal(c, importFileName)
			if err != nil {
				return err
			}

			c.imports = append(c.imports, ci)
		}
	}

	return nil
}

func (c *Config) start(r *DBRest) error {
	err := c.DB.start()
	if err != nil {
		return err
	}

	for _, handler := range c.Rest {
		handler.init(c.prefix, c.DB, r.rest)
		c.handlers = append(c.handlers, handler)
	}

	if len(c.imports) > 0 {
		for _, ci := range c.imports {
			err = ci.start(r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) stop() {
	if c.DB != nil {
		c.DB.stop()
	}

	if len(c.imports) > 0 {
		for _, ci := range c.imports {
			ci.stop()
		}
	}

}
