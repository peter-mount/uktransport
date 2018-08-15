package naptanimport

import(
  "encoding/csv"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
)

func (a *NaptanImport) railRef( r io.ReadCloser ) error {
  err := a.db.Update( func( tx *db.Tx ) error {
    stmt, err := tx.Prepare( "INSERT INTO naptan.rail VALUES ($1,$2,$3,$4,$5,$6,$7,$8,ST_SetSRID(ST_MakePoint($9,$10), 27700))" )
    if err != nil {
      return err
    }

    tx.OnCommitCluster( "naptan.rail", "rail_geom" )

    _, err = tx.DeleteFrom( "naptan.rail" )
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

        _, err := stmt.Exec(
          rec[0],
          rec[1],
          rec[2],
          rec[3],
          rec[8],
          rec[9],
          rec[10],
          rec[11],
          rec[6],
          rec[7],
        )
        if err != nil {
          return err
        }
        ic++

      }
    }

    log.Println( "Inserted", ic )

    return nil
  } )

  return err
}
