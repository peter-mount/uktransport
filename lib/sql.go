package lib

import (
//  "database/sql"
  "bufio"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
  "strings"
)

type SqlService struct {
  install      *bool
  // The DB
  db           *db.DBService
  // Optional, schema that if absent causes an install on startup
  Schema        string
}

func (a *SqlService) Name() string {
  return "SqlService"
}

func (a *SqlService) Init( k *kernel.Kernel ) error {
  a.install = flag.Bool( "install-db", false, "Install/Reinstall the database schema. THIS WILL ERASE ANY EXISTING DATA" )

  dbservice, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }
  a.db = (dbservice).(*db.DBService)

  return nil
}

func (a *SqlService) Start() error {

  if !*a.install && a.Schema != "" {
    // Check schema name exists & install if it doesn't
    exists, err := a.SchemaExists( a.Schema )
    if err != nil {
      return err
    }
    *a.install = !exists
  }

  if *a.install {
    err := a.Install()
    if err != nil {
      return err
    }
  }

  return nil
}

func (a *SqlService) Installed() bool {
  if a.install == nil {
    return false
  }
  return *a.install
}

func (a *SqlService) SchemaExists( schema string ) (bool, error) {
  row := a.db.QueryRow( "SELECT exists(select schema_name FROM information_schema.schemata WHERE schema_name = $1)", schema )
  var exists bool
  err := row.Scan( &exists )
  if err != nil {
    return false, err
  }
  return exists, nil
}

func (a *SqlService) Install() error {
  for _, name := range AssetNames() {
    log.Println( "Executing", name )

    asset, err := AssetString( name )
    if err != nil {
      return err
    }

    sr := strings.NewReader( asset )
    err = a.importSql( name, sr )
    if err != nil {
      return err
    }
  }
  return nil
}

func (a *SqlService) importSql( fn string, r *strings.Reader ) error {
  br := bufio.NewReader( r )
  sql := ""
  lc := 0
  for {
    s, err := br.ReadString( '\n' )
    if err != nil {
      if err == io.EOF {
        break
      }
      return err
    }

    if len(s)<3 || !(s[0] == s[1] && s[0] == '-') {
      sql = sql + s
    }

    if len(sql) > 2 && sql[len(sql)-2] == ';' {
      err := a.exec( fn, lc, sql )
      if err != nil {
        return err
      }
      sql = ""
    }

    lc++
  }

  if sql != "" {
    err := a.exec( fn, lc, sql )
    if err != nil {
      return err
    }
  }

  return nil
}

func (a *SqlService) exec( fn string, lc int, sql string ) error {
  _, err := a.db.Exec( sql )
  if err != nil {
    log.Println( fn, ":", lc, sql )
    return err
  }
  return nil
}
