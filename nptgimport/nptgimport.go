package nptgimport

import (
  //"database/sql"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
  "github.com/peter-mount/uktransport/lib"
)

type NptgImport struct {
  retrieve     *bool
  importdata   *bool

  dbdir        *string

  // The DB
  db           *db.DBService
  sql          *lib.SqlService
  zipImporter  *lib.ZipImporter
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
  a.db = (dbservice).(*db.DBService)

  sqlservice, err := k.AddService( &lib.SqlService{} )
  if err != nil {
    return err
  }
  a.sql = (sqlservice).(*lib.SqlService)
  a.sql.Schema = "nptg"

  zipImporter, err := k.AddService( lib.NewZipImporter(
    a.zipFile(),
    lib.ZipImportHandlerMap{
      "PlusbusMapping.csv": a.plusBusMapping,
    } ) )
  if err != nil {
    return err
  }
  a.zipImporter = (zipImporter).(*lib.ZipImporter)

  return nil
}

func (a *NptgImport) PostInit() error {
  if *a.dbdir == "" {
    *a.dbdir = "/database"
  }
  a.zipImporter.SetDir( *a.dbdir )

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

  // A retrieve, forced import or the schema being Installed then import the zip
  if *a.retrieve || *a.importdata || a.sql.Installed() {
    err := a.zipImporter.Import()
    if err != nil {
      return err
    }
  }

  return nil
}
