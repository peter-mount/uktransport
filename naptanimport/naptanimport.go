package naptanimport

import (
//  "database/sql"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
  "github.com/peter-mount/uktransport/lib"
)

type NaptanImport struct {
  retrieve     *bool
  importdata   *bool

  dbdir        *string

  // The DB
  db           *db.DBService
  sql          *lib.SqlService
}

func (a *NaptanImport) Name() string {
  return "nptgimport"
}

func (a *NaptanImport) Init( k *kernel.Kernel ) error {
  a.retrieve = flag.Bool( "retrieve", false, "Retrieve latest data" )
  a.importdata = flag.Bool( "import", false, "Import latest data" )
  a.dbdir = flag.String( "d", "", "Directory to store files, defaults to /database" )

  dbservice, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }
  a.db = (dbservice).(*db.DBService)

  sqlservice, err := k.AddService( &lib.SqlService{} )
  if err != nil {
    return err
  }
  a.sql = (sqlservice).(*lib.SqlService)

  return nil
}

func (a *NaptanImport) PostInit() error {
  if *a.dbdir == "" {
    *a.dbdir = "/database"
  }

  return nil
}

func (a *NaptanImport) Run() error {
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

  // A retrieve, forced import or the schema being Installed then import the zip
  if *a.retrieve || *a.importdata || a.sql.Installed() {
    zipImporter := lib.NewZipImporter( lib.ZipImportHandlerMap{
      "RailReferences.csv": a.railRef,
      "Stops.csv": a.stops,
    } )

    err := zipImporter.ImportZipFile( a.zipFile() )
    if err != nil {
      return err
    }
  }

  return nil
}
