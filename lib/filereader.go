package lib

import (
  "errors"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "io/ioutil"
  "path/filepath"
  "os"
  "log"
)

type FileReader struct {
  // DOS CR/LF record separators
  dos              *bool
  // RS record separator
  rs               *bool
  // STX/ETX record separator
  stxetx           *bool
  // Unix LF record separator
  unix             *bool

  pattern          *string

  parser            FileParser
  recordProcessor   RecordProcessor
}

func (a *FileReader) Name() string {
  return "fileReader"
}

func (a *FileReader) Init( k *kernel.Kernel ) error {
  a.dos = flag.Bool( "dos", false, "Use DOS CR/LF delimited records" )
  a.rs = flag.Bool( "rs", false, "Use RS delimited records" )
  a.stxetx = flag.Bool( "stxetx", false, "Use STX/ETX delimited records" )
  a.unix = flag.Bool( "unix", false, "Use Unix LF delimited records" )

  a.pattern = flag.String( "pattern", "", "Pattern to filter files" )

  return nil
}

func (a *FileReader) PostInit() error {
  // Setup the record separators, only one is allowed
  validate := func(a, b, c bool) error {
    if a || b || c {
      return errors.New( "Only one of -dos, -rs, -stxetx or -unix is permitted" )
    }
    return nil
  }
  if *a.dos {
    if err := validate( *a.rs, *a.stxetx, *a.unix ); err != nil {
      return err
    }
    a.parser = &DosParser{}
  } else if *a.rs {
    if err := validate( *a.dos, *a.stxetx, *a.unix ); err != nil {
      return err
    }
    a.parser = &UnixParser{ Char: 30 }
  } else if *a.stxetx {
    if err := validate( *a.dos, *a.rs, *a.unix ); err != nil {
      return err
    }
    a.parser = &StxEtxParser{}
  } else if *a.unix {
    if err := validate( *a.dos, *a.rs, *a.stxetx ); err != nil {
      return err
    }
    a.parser = &UnixParser{ Char: 10 }
  } else {
    return errors.New( "One of -dos, -rs, -stxetx or -unix is required" )
  }

  a.recordProcessor = func(b []byte ) error {
    s := string(b[:])
    log.Println( s )
    return nil
  }

  return nil
}

func (a *FileReader) ReadFile( s string ) error {
  f, err := os.Open( s )
  if err != nil {
    return err
  }
  defer f.Close()

  fi, err := f.Stat()
  if err != nil {
    return err
  }

  // If arg is a directory then walk it
  if fi.IsDir() {
    return filepath.Walk( s,
      func(path string, info os.FileInfo, err error) error {
        if err != nil {
          return err
        }

        if info.IsDir() {
          return nil
        }

        if *a.pattern != "" {
          m, err := filepath.Match( *a.pattern, filepath.Base( path ) )
          if !m || err != nil {
            return err
          }
        }

        return a.ReadFile( path )
      } )
  }

  // Read the file
  b, err := ioutil.ReadAll( f )
  if err != nil {
    return err
  }

  return a.parser.Parse( b, a.recordProcessor )
}
