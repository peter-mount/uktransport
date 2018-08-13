package nptgimport

import(
  "database/sql"
  "encoding/csv"
  "io"
  "log"
  "strings"
)

// plusBusMapping handles the import of the PlusBusMapping.csv file.
// This file contains the polygon data, one point per record for the area
// defining the plusBus area
func (a *NptgImport) plusBusMapping( r io.ReadCloser ) error {
  err := a.Update( func( tx *sql.Tx ) error {
    result, err := tx.Exec( "DELETE FROM nptg.plusbus" )
    if err != nil {
      return err
    }
    ra, err := result.RowsAffected()
    if err != nil {
      return err
    }
    log.Println( "Deleted", ra )

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
            err := plusBusMapping_persist( tx, crec, coords )
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
      err := plusBusMapping_persist( tx, crec, coords )
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

func plusBusMapping_persist( tx *sql.Tx, rec []string, coords []string ) error {
  // Ensure we are closed
  if coords[len(coords)-1] != coords[0] {
    coords = append( coords, coords[0] )
  }

  sql := "INSERT INTO nptg.plusbus VALUES ('" +
    rec[0] + "'," +
    "NULL," +
    "NULL,'" +
    rec[5] + "','" +
    rec[6] + "','" +
    rec[7] + "','" +
    rec[8] + "'," +
    "ST_MakePolygon(ST_GeomFromText('LINESTRING(" + strings.Join( coords, "," ) + ")', 27700))" +
    ")"

  _, err := tx.Exec( sql )
  if err != nil {
    return err
  }

  return nil
}
