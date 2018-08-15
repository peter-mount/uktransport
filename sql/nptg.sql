-- ================================================================================
-- nptg schema
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS nptg;

-- ================================================================================
-- plusbus consists of a polygon for each Plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS nptg.plusbus CASCADE;

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
DROP TABLE IF EXISTS nptg.plusbuszone CASCADE;

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

-- ================================================================================
-- railplusbus is a view of naptan.rail that is contained within a plusbus zone
-- Specifically this means rail stations within each zone.
-- This view inherits the geometry from naptan.rail so can be used as a Point feature
-- ================================================================================
CREATE VIEW nptg.railplusbus
  AS SELECT z.zone AS plusbuszone, r.*
    FROM naptan.rail r
      INNER JOIN nptg.plusbus z ON ST_Contains( z.geom, r.geom );
