package nptgimport

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
  "time"
)

func (a *NptgImport) zipFile() string {
  return *a.dbdir + "/Nptgcsv.zip"
}

func (a *NptgImport) retrieveRequired() (bool, error) {
  fi, err := os.Stat( a.zipFile() )
  if err != nil {
    if os.IsNotExist( err ) {
      return true, nil
    }
    return false, err
  }

  n := time.Now()

  return n.Sub( fi.ModTime() ) > (12 * time.Hour), nil
}

func (a *NptgImport) retrieveData() error {
  log.Println( "Retrieving nptg data" )

  url := "http://naptan.app.dft.gov.uk/datarequest/nptg.ashx?format=csv"
  req, err := http.NewRequest( "GET", url, nil )
  if err != nil {
    return err
  }

  log.Println( "Retrieving", url )
  resp, err := http.DefaultClient.Do( req )
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  if resp.StatusCode < 200 || resp.StatusCode >= 300 {
    return fmt.Errorf( "Request returned %d: %s", resp.StatusCode, resp.Status )
  }
  log.Printf( "Retrieved %d bytes", resp.ContentLength )

  // Copy the body to a temporary file which is deleted when retrieve exits
  file, err := os.OpenFile( a.zipFile(), os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644 )
  if err != nil {
    return err
  }
  defer file.Close()

  _, err = io.Copy( file, resp.Body )
  if err != nil {
    return err
  }

  return nil
}
