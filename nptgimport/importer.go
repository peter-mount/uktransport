package nptgimport

import (
  "archive/zip"
  "io"
  "log"
)

func (a *NptgImport) importData() error {

  lookup := map[string]func(io.ReadCloser) error {
    "PlusbusMapping.csv": a.plusBusMapping,
  }

  r, err := zip.OpenReader( a.zipFile() )
  if err != nil {
      return err
  }
  defer r.Close()

  for _, f := range r.File {
    if fh, ok := lookup[f.Name]; ok {
      err := a.importFile( f, fh )
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (a *NptgImport) importFile( f *zip.File, h func(io.ReadCloser) error ) error {
  log.Println( "Importing", f.Name )
  rc, err := f.Open()
  if err != nil {
    return err
  }
  defer rc.Close()
  
  return h(rc)
}
