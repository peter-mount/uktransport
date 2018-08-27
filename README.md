# uktransport

This docker image contains a suite of command line utilities that support the retrieval of various UK Transport OpenData
and the importing of that data into a PostGIS database.

We use PostGIS rather than plain PostgreSQL as some of this data is geographic in nature.

Instructions on how to use this image will appear in the wiki.

## Available commands

### cif
#### cifimport
cifimport handles the importing of the Network Rail CIF timetable, providing both schedules and tiploc
entries describing the UK Rail Network. This is from the [nrod-cif](https://github.com/peter-mount/nrod-cif) project.

#### cifretrieve
cifretrieve handles the retrieval Network Rail CIF timetable in CF format, providing both schedules and tiploc
entries describing the UK Rail Network. This is from the [nrod-cif](https://github.com/peter-mount/nrod-cif) project.

### NaPTAN
#### naptanimport
naptanimport retrieves and imports the NaPTAN dataset directly from the UK's Department of Transport.
It contains details about the locations of Airports, Railway stations, Bus stops for the entire country.

#### nptgimport
nptgimport retrieves and imports the NPTG dataset directly from the UK's Department of Transport.
This dataset contains details about localities within the UK, for example where a specific town is located.
It also includes geographic coverage of the  PlusBus zones (a type of Bus ticket valid with Rail tickets).
