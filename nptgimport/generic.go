package nptgimport

import(
  "database/sql"
  "encoding/csv"
  "fmt"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
  "strings"
)

// genericImport imports the csv file into a table with the file name
func (a *NptgImport) genericImport( n string, r io.ReadCloser ) error {
  err := a.db.Update( func( tx *db.Tx ) error {

    var stmt *sql.Stmt

    table := "nptg." + n[:len(n)-4]

    tx.OnCommitVacuumFull( table )

    _, err := tx.DeleteFrom( table )
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

      if lc == 0 {
        // Header line so create a statement with the same number of columns
        s := []string{"INSERT INTO " + table + " VALUES ("}
        f := "$%d"
        for i, _ := range rec {
          s = append( s, fmt.Sprintf( f, i+1 ) )
          f = ",$%d"
        }
        s = append( s, ")" )
        stmt, err = tx.Prepare( strings.Join( s, "" ) )
        if err != nil {
          return err
        }
        defer stmt.Close()
      } else {
        // Import the line
        args := make( []interface{}, len(rec) )
        for i, v := range rec {
          if v == "" {
            args[i] = nil
          } else {
            args[i] = v
          }
        }
        _, err = stmt.Exec( args... )
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
