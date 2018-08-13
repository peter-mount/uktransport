package naptanimport

import(
  "database/sql"
  "encoding/csv"
  "io"
  "log"
)

func nullIfEmpty( s string ) *string {
  if s == "" {
    return nil
  }
  return &s
}

func (a *NaptanImport) stops( r io.ReadCloser ) error {
  err := a.Update( func( tx *sql.Tx ) error {
    // Needed as naptan is not in UTF-8
    _, err := tx.Exec( "SET CLIENT_ENCODING TO WIN1252" )
    if err != nil {
      return err
    }

    result, err := tx.Exec( "DELETE FROM naptan.stops" )
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

        _, err := tx.Exec( "INSERT INTO naptan.stops VALUES (" +
          " $1, $2, $3, $4, $5, $6, $7, $8, $9,$10," +
          "$11,$12,$13,$14,$15,$16,$17,$18,$19,$20," +
          "$21,$22,$23,$24,$25,$26,$27,$28,$29,$30," +
          "$31,$32," +
           "ST_SetSRID(ST_MakePoint($33,$34), 27700))",
          rec[0],
          rec[1],
          rec[2],
          rec[3],
          rec[4],
          rec[6],
          rec[8],
          rec[10],
          rec[12],
          rec[14],

          rec[16], // bearing
          rec[17],
          rec[18],
          rec[19],
          rec[20],
          rec[21], // town
          rec[23],
          rec[27], // east
          rec[28],
          rec[29], // long

          rec[30], // lat
          rec[31],
          rec[32],
          rec[33],
          rec[34],
          rec[35], // notes
          rec[37],
          nullIfEmpty(rec[38]), // created
          nullIfEmpty(rec[39]),
          rec[40],

          rec[41],
          rec[42],

          rec[27], // east
          rec[28],
        )
        if err != nil {
          log.Println( "Failed atco", rec[0] )
          return err
        }
        ic++

      }
    }

    log.Println( "Inserted", ic )

    return nil
  } )
  if err != nil {
    return err
  }

  return nil
}
