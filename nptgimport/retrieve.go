package nptgimport

import (
  "github.com/peter-mount/uktransport/lib"
)

func (a *NptgImport) zipFile() string {
  return *a.dbdir + "/Nptgcsv.zip"
}

func (a *NptgImport) retrieveRequired() (bool, error) {
  return lib.RetrieveRequired( a.zipFile() )
}

func (a *NptgImport) retrieveData() error {
  return lib.RetrieveHttp( a.zipFile(), "http://naptan.app.dft.gov.uk/datarequest/nptg.ashx?format=csv" )
}
