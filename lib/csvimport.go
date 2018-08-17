package lib

import(
  "database/sql"
  "encoding/csv"
  "fmt"
  "github.com/peter-mount/golib/kernel/db"
  "io"
  "log"
  "strings"
)

// csvimport covers the entire import
type csvimport struct {
  // The transaction
  tx             *db.Tx
  // The table name
  table           string
  // Prepared statement
  stmt           *sql.Stmt
  // map of columns by name
  names           map[string]*column
  // slice of columns as seen in import
  columns      []*column
  // Number of lines imported
  importCount     int
  // true if we have geometry
  hasGeometry     bool
  geom_srid       int
  geom_east       string
  geom_north      string
  geom_index      string
}

// Internal column mapping column names to sql fields
type column struct {
  // Name of column
  name    string
  // Index in csv (0...)
  column  int
}

const (
  easting = "Easting"
  northing = "Northing"
  longitude = "Longitude"
  latitude = "Latitude"
)

func (c *csvimport) close() {
  if c.stmt != nil {
    c.stmt.Close()
  }
}

func (c *csvimport) columnExists( n string ) bool {
  _, exists := c.names[ n ]
  return exists
}

func (c *csvimport) checkGeometry() bool {
  if c.columnExists( easting ) && c.columnExists( northing ) {
    c.geom_srid = 27700
    c.geom_east = easting
    c.geom_north = northing
  } else if c.columnExists( longitude ) && c.columnExists( latitude ) {
    c.geom_srid = 4326
    c.geom_east = longitude
    c.geom_north = latitude
  }
  c.hasGeometry = c.geom_srid != 0 && c.geom_east != "" && c.geom_north != ""
  return c.hasGeometry
}

// Generate the import statement
func (c *csvimport) generateStatement( row []string ) error {
  // Header line so create a statement with the same number of columns
  s := []string{"INSERT INTO " + c.table + " ("}

  for i, col := range c.columns {
    if i > 0 {
      s = append( s, "," )
    }
    s = append( s, col.name )
  }

  s = append( s, ") VALUES (" )

  f := "$%d"
  for i, _ := range c.columns {
    s = append( s, fmt.Sprintf( f, i+1 ) )
    f = ",$%d"
  }

  s = append( s, ")" )

  log.Println( strings.Join( s, "" ) )

  stmt, err := c.tx.Prepare( strings.Join( s, "" ) )
  c.stmt = stmt
  return err
}

func (c *csvimport) parseRow( row []string ) error {
  if c.stmt == nil {
    return c.parseHeader( row )
  }
  return c.insertRow( row )
}

// Parse the first header row
func (c *csvimport) parseHeader( row []string ) error {
  c.names = make( map[string]*column )
  for i, n := range row {
    if _, exists := c.names[ n ]; exists {
      return fmt.Errorf( "Unable to import csv, duplicate column header \"%s\"", n )
    }

    col := &column{ name: n, column: i }
    c.names[ n ] = col
    c.columns = append( c.columns, col )
  }

  // create the Stmt
  err := c.generateStatement( row )
  if err != nil {
    log.Println( "Failed to prepare statement" )
    return err
  }

  // Check to see if we have geometry
  if c.checkGeometry() {
    // Update the geometry prior to the commit
    c.tx.BeforeCommit( c.updateGeometry )

    // Cluster on the geometry at the end
    c.tx.OnCommitCluster( c.table, c.geom_index )
  } else {
    // Vacuum at the end when we have no geometry
    c.tx.OnCommitVacuumFull( c.table )
  }

  // Now truncate the table
  _, err = c.tx.DeleteFrom( c.table )
  return err
}

// isNull returns nil if v is "" or contains 0x00 otherwise v
func isNull( v string ) interface{} {
  if v == "" || v == "\x00" {
    return nil
  }
  return v
}

func (c *csvimport) insertRow( row []string ) error {
  var args []interface{}
  for _, col := range c.columns {
    args = append( args, isNull( row[col.column] ) )
  }

  _, err := c.stmt.Exec( args... )
  if err != nil {
    log.Println( "Insert failed:", c.importCount+1,"\n\"", strings.Join(row,"\",\""), "\"" )
    return err
  }

  c.importCount++
  return nil
}

// genericImport imports the csv file into a table with the file name
func (c *csvimport) updateGeometry( tx *db.Tx ) error {

  if c.hasGeometry {
    log.Printf(
      "Updating geometry using %s(%s,%s) EPSG:%d",
      c.table,
      c.geom_east,
      c.geom_north,
      c.geom_srid )

    result, err := c.tx.Exec( fmt.Sprintf(
      "UPDATE %s SET geom = ST_SetSRID(ST_MakePoint(%s,%s),%d) WHERE %s IS NOT NULL AND %s IS NOT NULL",
      c.table,
      c.geom_east, c.geom_north,
      c.geom_srid,
      c.geom_east, c.geom_north,
    ) )
    if err != nil {
      return err
    }

    ra, err := result.RowsAffected()
    if err != nil {
      return err
    }

    log.Printf( "Updated %d entries", ra )
  }

  return nil
}

// genericImport imports the csv file into a table with the file name
func (a *SqlService) CSVImport( n string, r io.ReadCloser ) error {
  if a.Schema == "" {
    return fmt.Errorf( "No schema defined for CSVImport: %s", n)
  }

  tableName := n[:len(n)-4]
  state := csvimport{
    table: a.Schema + "." + tableName,
    geom_index: tableName + "_geom",
  }
  defer state.close()

  return a.db.Update( func( tx *db.Tx ) error {
    state.tx = tx

    // Needed as csv files are not in UTF-8
    err := tx.SetEncodingWIN1252()
    if err != nil {
      return err
    }

    rdr := csv.NewReader( r )
    for {
      rec, err := rdr.Read()
      if err == io.EOF {
        break
      }
      if err != nil {
        return err
      }

      err = state.parseRow( rec )
      if err != nil {
        return err
      }
    }

    log.Println( "Inserted", state.importCount )

    return nil
  } )
}
