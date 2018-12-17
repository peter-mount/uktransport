# uktransport

This docker image contains a suite of command line utilities that support the retrieval of various UK Transport OpenData.

Some of these utilities will write the retrieved data into a PostGIS database whilst others pass that data to a local RabbitMQ server.

We use PostGIS rather than plain PostgreSQL as some of this data is geographic in nature.

Instructions on how to use this image will appear in the wiki.

Currently this image supports the amd64, arm64v8 & arm32v7 (a.k.a. Raspberry PI 3B, 3B+ & 3A) architectures.

## Data retrieval commands

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

## Data manipulation commands

### dataretriever

dataretriever is a generic tool to retrieve data and pass it to a RabbitMQ instance.
This is from the [dataretriever](https://github.com/peter-mount/dataretriever) project.

It currently has two modes of operation:
* Retrieve via http/https at regular intervals data and submit the response as a message.
* Connect to a remote message broker using Stomp and submit messages to RabbitMQ. For Rail open data this suppots the NROD feed from Network Rail but *not* the Darwin Push Port feed.

### publishmq

publishmq is a utility currently being written (so not yet usable) to parse archived logs taken from the open data feeds and resubmit them to a RabbitMQ instance. It's mainly for use in testing the code that parses the data feeds.

### rabtap

rabtap is a Swiss army knife for RabbitMQ. Tap/Pub/Sub messages, create/delete/bind queues and exchanges, inspect broker.

Full documentation is in it's repository: https://github.com/jandelgado/rabtap
