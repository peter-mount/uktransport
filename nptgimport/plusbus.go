package nptgimport

import(
  "database/sql"
  "encoding/csv"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
  "strings"
)

// plusBusMapping handles the import of the PlusBusMapping.csv file.
// This file contains the polygon data, one point per record for the area
// defining the plusBus area
func (a *NptgImport) plusBusMapping( n string, r io.ReadCloser ) error {
  err := a.db.Update( func( tx *db.Tx ) error {

    // Needed as csv files are not in UTF-8
    err := tx.SetEncodingWIN1252()
    if err != nil {
      return err
    }

    table := "nptg.PlusbusMapping"
    geom := "PlusbusMapping_geom"

    stmt, err := tx.Prepare( "INSERT INTO " + table + " VALUES ($1,$2,$3,$4,$5,ST_MakePolygon(ST_GeomFromText($6, 27700)))" )
    if err != nil {
      return err
    }
    defer stmt.Close()

    tx.OnCommitCluster( table, geom )

    _, err = tx.DeleteFrom( table )
    if err != nil {
      return err
    }

    lc := 0
    ic := 0
    code := ""
    var crec []string
    var coords []string
    rdr := csv.NewReader( r )
    for {
      rec, err := rdr.Read()
      if err == io.EOF {
        break
      }
      if err != nil {
        return err
      }

      lc++
      if lc > 1 {
        if code != rec[0] {
          if code != "" {
            err := plusBusMapping_persist( stmt, crec, coords )
            if err != nil {
              return err
            }
            ic++
          }

          coords = nil
          crec = rec
          code = rec[0]
        }
        coords = append( coords, rec[3] + " " + rec[4] )
      }
    }

    // Record the last entry
    if code != "" {
      err := plusBusMapping_persist( stmt, crec, coords )
      if err != nil {
        return err
      }
      ic++
    }

    log.Println( "Inserted", ic )

    return nil
  } )
  if err != nil {
    return err
  }

  return nil
}

func plusBusMapping_persist( stmt *sql.Stmt, rec []string, coords []string ) error {
  // Ensure we are closed
  if coords[len(coords)-1] != coords[0] {
    coords = append( coords, coords[0] )
  }

  _, err := stmt.Exec(
    rec[0],
    rec[5],
    rec[6],
    rec[7],
    rec[8],
    "LINESTRING(" + strings.Join( coords, "," ) + ")",
  )

  return err
}
