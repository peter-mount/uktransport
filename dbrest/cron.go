package dbrest

import (
	"database/sql"
	"gopkg.in/robfig/cron.v2"
	"log"
)

// CronJob represents a postgresql function thats called at regular intervals.
// These type of calls cannot take any parameters so just the function name needs to be provided.
type CronJob struct {
	// The cron schedule
	Schedule string `yaml:"schedule"`
	// The Postgres function to call
	Function string `yaml:"function"`
	// Log any errors to the console
	LogErrors bool `yaml:"logErrors"`
	statement *sql.Stmt
}

func (j *CronJob) init(db *sql.DB, c *cron.Cron) error {
	stmt, err := db.Prepare("SELECT " + j.Function)
	if err != nil {
		return err
	}
	j.statement = stmt

	_, err = c.AddFunc(j.Schedule, j.execute)
	if err != nil {
		return err
	}

	return nil
}

func (j *CronJob) execute() {
	_, err := j.statement.Exec(j.Function)
	if err != nil && j.LogErrors {
		log.Println("Failed to invoke:", j.Function, err)
	}
}
