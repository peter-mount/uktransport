-- ================================================================================
-- naptan
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS naptan;

-- ================================================================================
-- railreference
-- ================================================================================
DROP TABLE IF EXISTS naptan.rail;

CREATE TABLE naptan.rail (
  atco          NAME NOT NULL,
  tiploc        NAME,
  crs           CHAR(3),
  name          NAME,
  created       TIMESTAMP WITHOUT TIME ZONE,
  modified      TIMESTAMP WITHOUT TIME ZONE,
  revision      INTEGER,
  modification  NAME,
  PRIMARY KEY (atco)
);

CREATE INDEX rail_t ON naptan.rail(tiploc);
CREATE INDEX rail_c ON naptan.rail(crs);
CREATE INDEX rail_n ON naptan.rail(lower(name));


-- geometry
SELECT addgeometrycolumn( '', 'naptan', 'rail', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX rail_geom ON naptan.rail USING GIST (geom);

-- ================================================================================
-- stops
-- ================================================================================
DROP TABLE IF EXISTS naptan.stops;

CREATE TABLE naptan.stops (
  atco                    NAME NOT NULL,
  naptan                  NAME,
  platecode               NAME,
  cleardowncode           NAME,
  commonName              NAME,
  shortcommonname         NAME,
  landmark                NAME,
  street                  NAME,
  crossing                NAME,
  indicator               NAME,

  bearing                 NAME,
  nptglocalitycode        NAME,
  localityName            NAME,
  parentLocalityName      NAME,
  grandParentLocalityName NAME,
  town                    NAME,
  suburb                  NAME,
  easting                 INTEGER,
  northing                INTEGER,
  longitude               REAL,

  latitude                REAL,
  stoptype                NAME,
  busstoptype             NAME,
  timingstatus            NAME,
  defaultwaittime         NAME,
  notes                   TEXT,
  adminareacode           NAME,
  created                 TIMESTAMP WITHOUT TIME ZONE,
  modified                TIMESTAMP WITHOUT TIME ZONE,
  revision                INTEGER,

  modification            NAME,
  status                  NAME,
  PRIMARY KEY (atco)
);

-- geometry
SELECT addgeometrycolumn( '', 'naptan', 'stops', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX stops_geom ON naptan.stops USING GIST (geom);
