package lib

import (
  "bytes"
  "unicode/utf8"
)

type RecordProcessor func([]byte) error
type FileParser interface {
  Parse( []byte, RecordProcessor ) error
}

type StxEtxParser struct {
  b     bytes.Buffer
  inRec bool
}
func (p *StxEtxParser) Parse( s []byte, f RecordProcessor ) error {
  p.b.Reset()
  p.inRec = false

  for len(s) > 0 {
    switch s[0] {
      case 2:
        p.b.Reset()
        p.inRec = true
        s = s[1:]

      case 3:
        if p.inRec && p.b.Len() > 0 {
          if err := f(p.b.Bytes()); err != nil {
            return err
          }
          p.b.Reset()
          p.inRec = false
        }
        s = s[1:]

      default:
        r, size := utf8.DecodeRune(s)

        if p.inRec {
          if _, err := p.b.WriteRune( r ); err != nil {
            return err
          }
        }

        s = s[size:]
    }
	}

  // Don't handle unterminated entries as it must be unterminated with Etx
  if p.b.Len() > 0 {
    if err := f(p.b.Bytes()); err != nil {
      return err
    }
  }

  return nil
}

// UnixParser handles parsing for a single char delimiter, usinally LF (10)
// hence it's name UnixParser, but can be used for delimiting against any
// other char, e.g. RS (30)
type UnixParser struct {
  Char  byte
  b     bytes.Buffer
}
func (p *UnixParser) Parse( s []byte, f RecordProcessor ) error {
  p.b.Reset()

  for len(s) > 0 {
    switch s[0] {
    case p.Char:
        if p.b.Len() > 0 {
          if err := f(p.b.Bytes()); err != nil {
            return err
          }
          p.b.Reset()
        }
        s = s[1:]

      default:
        r, size := utf8.DecodeRune(s)

        if _, err := p.b.WriteRune( r ); err != nil {
          return err
        }

        s = s[size:]
    }
	}

  // Handle unterminated entries
  if p.b.Len() > 0 {
    if err := f(p.b.Bytes()); err != nil {
      return err
    }
  }

  return nil
}

// DosParser handles parsing for CRLF delimited files
type DosParser struct {
  b     bytes.Buffer
}
func (p *DosParser) Parse( s []byte, f RecordProcessor ) error {
  p.b.Reset()

  for len(s) > 0 {
    if len(s) >1 && s[0] == 13 && s[1] == 10 {
      if p.b.Len() > 0 {
        if err := f(p.b.Bytes()); err != nil {
          return err
        }
        p.b.Reset()
      }
      s = s[2:]
    } else {
      r, size := utf8.DecodeRune(s)

      if _, err := p.b.WriteRune( r ); err != nil {
        return err
      }

      s = s[size:]
    }
  }

  // Handle unterminated entries
  if p.b.Len() > 0 {
    if err := f(p.b.Bytes()); err != nil {
      return err
    }
  }

  return nil
}
