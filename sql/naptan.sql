-- ================================================================================
-- naptan
-- ================================================================================

CREATE EXTENSION postgis;

CREATE SCHEMA naptan;

DROP TABLE naptan.rail;

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
CREATE INDEX plusbus_geom ON naptan.rail USING GIST (geom);
