-- ================================================================================
-- nptg schema
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS nptg;

-- ================================================================================
-- plusbus consists of a polygon for each Plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS nptg.plusbus;

CREATE TABLE nptg.plusbus (
  zone          NAME NOT NULL,
  created       TIMESTAMP WITHOUT TIME ZONE,
  modified      TIMESTAMP WITHOUT TIME ZONE,
  revision      INTEGER,
  modification  NAME,
  PRIMARY KEY (zone)
);

-- geometry
SELECT addgeometrycolumn( '', 'nptg', 'plusbus', 'geom', 27700, 'POLYGON', 2, true);
CREATE INDEX plusbus_geom ON nptg.plusbus USING GIST (geom);

-- ================================================================================
-- plusbuszone Holds the metadata for a Plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS nptg.plusbuszone;

CREATE TABLE nptg.plusbuszone (
  zone          NAME NOT NULL,
  name          NAME,
  namelang      NAME,
  country       NAME,
  created       TIMESTAMP WITHOUT TIME ZONE,
  modified      TIMESTAMP WITHOUT TIME ZONE,
  revision      INTEGER,
  modification  NAME,
  PRIMARY KEY (zone)
);
