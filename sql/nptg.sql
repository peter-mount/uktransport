-- ================================================================================
-- nptg schema
-- ================================================================================

CREATE EXTENSION postgis;

CREATE SCHEMA nptg;

DROP TABLE nptg.plusbus;

CREATE TABLE nptg.plusbus (
  zone          NAME NOT NULL,
  name          NAME,
  country       NAME,
  created       TIMESTAMP WITHOUT TIME ZONE,
  modified      TIMESTAMP WITHOUT TIME ZONE,
  revision      INTEGER,
  modification  NAME,
  PRIMARY KEY (zone)
);

-- geometry
SELECT addgeometrycolumn( '', 'nptg', 'plusbus', 'geom', 27700, 'POLYGON', 2, true);
CREATE INDEX plusbus_geom ON nptg.plusbus USING GIST (geom);
