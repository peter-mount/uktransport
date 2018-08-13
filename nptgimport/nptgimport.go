package nptgimport

import (
  "database/sql"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
)

type NptgImport struct {
  retrieve     *bool
  importdata   *bool

  dbdir        *string

  // The DB
  dbService    *db.DBService
}

func (a *NptgImport) Name() string {
  return "nptgimport"
}

func (a *NptgImport) Init( k *kernel.Kernel ) error {
  a.retrieve = flag.Bool( "retrieve", false, "Retrieve latest data" )
  a.importdata = flag.Bool( "import", false, "Import latest data" )
  a.dbdir = flag.String( "d", "", "Directory to store files, defaults to /database" )

  dbservice, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }
  a.dbService = (dbservice).(*db.DBService)

  return nil
}

func (a *NptgImport) PostInit() error {
  if *a.dbdir == "" {
    *a.dbdir = "/database"
  }

  return nil
}

func (a *NptgImport) Run() error {
  if !*a.retrieve {
    retr, err := a.retrieveRequired()
    if err != nil {
      return err
    }
    *a.retrieve = retr
  }

  if *a.retrieve {
    err := a.retrieveData()
    if err != nil {
      return err
    }
  }

  if *a.retrieve || *a.importdata {
    err := a.importData()
    if err != nil {
      return err
    }
  }

  return nil
}

func (c *NptgImport) Update( f func( *sql.Tx ) error ) error {
  tx, err := c.dbService.GetDB().Begin()
  if err != nil {
    return err
  }
  defer tx.Commit()

  err = f( tx )
  if err != nil {
    tx.Rollback()
    return err
  }

  return nil
}
