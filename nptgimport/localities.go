package nptgimport

import(
  "encoding/csv"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
)

// plusBusMapping handles the import of the PlusBusMapping.csv file.
// This file contains the polygon data, one point per record for the area
// defining the plusBus area
func (a *NptgImport) localities( n string, r io.ReadCloser ) error {
  table := "nptg.localities"

  return a.db.Update( func( tx *db.Tx ) error {

    // Needed as csv files are not in UTF-8
    err := tx.SetEncodingWIN1252()
    if err != nil {
      return err
    }

    stmt, err := tx.Prepare( "INSERT INTO " + table + " VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,ST_SetSRID(ST_MakePoint($17,$18), 27700))" )
    if err != nil {
      return err
    }
    defer stmt.Close()

    tx.OnCommitCluster( table, "localities_geom" )

    _, err = tx.DeleteFrom( table )
    if err != nil {
      return err
    }

    lc := 0
    ic := 0
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
        ic++
        _, err = stmt.Exec(
          rec[0],
          rec[1],
          rec[2],
          rec[3],
          rec[4],
          rec[5],
          rec[6],
          rec[7],
          rec[8],
          rec[9],
          rec[10],
          rec[11],
          rec[15],
          rec[16],
          rec[17],
          rec[18],
          rec[13], // easting
          rec[14], // northing
        )
        if err != nil {
          return err
        }
      }
    }

    log.Println( "Inserted", ic )

    return nil
  } )
}
