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
func (a *NptgImport) plusBusZone( r io.ReadCloser ) error {
  err := a.db.Update( func( tx *db.Tx ) error {

    stmt, err := tx.Prepare( "INSERT INTO nptg.plusbuszone VALUES ($1,$2,$3,$4,$5,$6,$7,$8)" )
    if err != nil {
      return err
    }
    defer stmt.Close()

    tx.OnCommitVacuumFull( "nptg.plusbuszone" )

    _, err = tx.DeleteFrom( "nptg.plusbuszone" )
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

      if lc > 0 {
        _, err = stmt.Exec(
          rec[0], rec[1], rec[2], rec[3],
          rec[4], rec[5], rec[6], rec[7],
        )
        if err != nil {
          return err
        }
        ic++
      }
      lc++
    }

    log.Println( "Inserted", ic )

    return nil
  } )
  if err != nil {
    return err
  }

  return nil
}
