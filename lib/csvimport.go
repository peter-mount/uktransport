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
  // trie if we have geometry, osgb first then wgs84 second
  hasGeometry     bool
  // true if OSGB geometry exists, i.e. Has Easting, Northing columns in ESRI:27700
  osgb            bool
  // true if WGS84 geometry exists, i.e. Has Longitude, Latitude columns in ESRI:4326
  wgs84           bool
  // The coordinate columns
  geom_east       int
  geom_north      int
  // The index for the geometry
  geomIndex       string
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

func (c *csvimport) checkGeometry() {
  if c.columnExists( easting ) && c.columnExists( northing ) {
    c.osgb = true
    c.geom_east = c.names[easting].column
    c.geom_north = c.names[northing].column
  } else if c.columnExists( longitude ) && c.columnExists( latitude ) {
    c.wgs84 = true
    c.geom_east = c.names[longitude].column
    c.geom_north = c.names[latitude].column
  }
  c.hasGeometry = c.osgb || c.wgs84
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

  if c.hasGeometry {
    s = append( s, ",geom" )
  }

  s = append( s, ") VALUES (" )

  f := "$%d"
  for i, _ := range c.columns {
    s = append( s, fmt.Sprintf( f, i+1 ) )
    f = ",$%d"
  }

  if c.hasGeometry {
    var srid int
    if c.osgb {
      srid = 27700
    } else if c.wgs84 {
      srid = 4326
    } else {
      return fmt.Errorf( "hasGeometry set but not osgb/wgs84" )
    }
    s = append( s, fmt.Sprintf( ",ST_SetSRID(ST_MakePoint($%d,$%d),%d)", c.geom_east, c.geom_north, srid ) )
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

  // Check to see if we have geometry
  c.checkGeometry()

  // create the Stmt
  err := c.generateStatement( row )
  if err != nil {
    return err
  }

  if c.hasGeometry {
    // Cluster on the geometry at the end
    c.tx.OnCommitCluster( c.table, c.geomIndex )
  } else {
    // Vacuum at the end when we have no geometry
    c.tx.OnCommitVacuumFull( c.table )
  }

  // Now truncate the table
  _, err = c.tx.DeleteFrom( c.table )
  return err
}

func isNull( v string ) interface{} {
  if v == "" {
    return nil
  }
  return v
}

func (c *csvimport) insertRow( row []string ) error {
  var args []interface{}
  for _, col := range c.columns {
    args = append( args, isNull( row[col.column] ) )
  }

  if c.hasGeometry {
    args = append( args, row[c.geom_east], row[c.geom_north] )
  }

  _, err := c.stmt.Exec( args... )
  if err != nil {
    return err
  }

  c.importCount++
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
    geomIndex: tableName + "_geom",
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
