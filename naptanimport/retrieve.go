package naptanimport

import (
  "github.com/peter-mount/uktransport/lib"
)

func (a *NaptanImport) zipFile() string {
  return *a.dbdir + "/Naptancsv.zip"
}

func (a *NaptanImport) retrieveRequired() (bool, error) {
  return lib.RetrieveRequired( a.zipFile() )
}

func (a *NaptanImport) retrieveData() error {
  return lib.RetrieveHttp( a.zipFile(), "http://naptan.app.dft.gov.uk/DataRequest/Naptan.ashx?format=csv" )
}
