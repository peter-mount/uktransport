package naptanimport

import(
  "encoding/csv"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
)

func (a *NaptanImport) stopPlusbusZones( r io.ReadCloser ) error {
  err := a.db.Update( func( tx *db.Tx ) error {
    stmt, err := tx.Prepare( "INSERT INTO naptan.stopplusbuszones VALUES ($1,$2)" )
    if err != nil {
      return err
    }

    // Cluster entries by plusbuz zone
    tx.OnCommitCluster( "naptan.stopplusbuszones", "stopplusbuszones_zone" )

    _, err = tx.DeleteFrom( "naptan.stopplusbuszones" )
    if err != nil {
      return err
    }

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

      _, err = stmt.Exec( rec[0], rec[1] )
      if err != nil {
        return err
      }
      ic++
    }

    log.Println( "Inserted", ic )

    return nil
  } )

  return err
}
