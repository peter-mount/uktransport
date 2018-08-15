package lib

import (
  "archive/zip"
  "io"
  "log"
)

// ZipImportHandler handles the import of a file inside the zip file
type ZipImportHandler func( io.ReadCloser ) error

type ZipImportHandlerMap map[string]ZipImportHandler

// ZipImporter handles the import of a zipFile
type ZipImporter struct {
  handlers ZipImportHandlerMap
}

func NewZipImporter( handlers ZipImportHandlerMap ) *ZipImporter {
  return &ZipImporter{ handlers: handlers }
}

// ImportZipFile scans a zip file and if an entry in a map exists for the filename
// will pass that to the handler
func (z *ZipImporter) ImportZipFile( fileName string ) error {

  r, err := zip.OpenReader( fileName )
  if err != nil {
      return err
  }
  defer r.Close()

  return z.ImportZipReader( r )
}

// ImportZipReader scans a zip file and if an entry in a map exists for the filename
// will pass that to the handler
func (z *ZipImporter) ImportZipReader( r *zip.ReadCloser ) error {

  for _, f := range r.File {
    if fh, ok := z.handlers[f.Name]; ok {
      err := z.importZipFile( f, fh )
      if err != nil {
        return err
      }
    }
  }

  return nil
}

func (z *ZipImporter) importZipFile( f *zip.File, h ZipImportHandler ) error {
  log.Println( "Importing", f.Name )
  rc, err := f.Open()
  if err != nil {
    return err
  }
  defer rc.Close()

  return h(rc)
}
