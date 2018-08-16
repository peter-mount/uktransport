package lib

import (
  "archive/zip"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "io"
  "log"
)

// ZipImportHandler handles the import of a file inside the zip file
type ZipImportHandler func( string, io.ReadCloser ) error

type ZipImportHandlerMap map[string]ZipImportHandler

// ZipImporter handles the import of a zipFile
type ZipImporter struct {
  name      string
  dir       string
  handlers  ZipImportHandlerMap
  zipfile  *string
}

func NewZipImporter( name string, handlers ZipImportHandlerMap ) *ZipImporter {
  return &ZipImporter{ name: name, handlers: handlers }
}

func (a *ZipImporter) Name() string {
  return "ZipImporter:" + a.name
}

func (a *ZipImporter) Init( k *kernel.Kernel ) error {
  a.zipfile = flag.String( "import-zip-file", "", "Import just one file within the source zip file" )
  return nil
}

func (z *ZipImporter) SetDir( dir string ) {
  z.dir = dir
}

// ImportZipFile scans a zip file and if an entry in a map exists for the filename
// will pass that to the handler
func (z *ZipImporter) Import() error {

  r, err := zip.OpenReader( z.dir + z.name )
  if err != nil {
      return err
  }
  defer r.Close()

  for _, f := range r.File {
    if fh, ok := z.handlers[f.Name]; ok {
      // If no filter or asking for a specific one
      if *z.zipfile == "" || *z.zipfile == f.Name {
        err := z.importZipFile( f, fh )
        if err != nil {
          return err
        }
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

  return h( f.Name, rc )
}
