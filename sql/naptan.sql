-- ================================================================================
-- naptan
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS naptan;

-- ================================================================================
-- AirReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.AirReferences CASCADE;

CREATE TABLE naptan.AirReferences (
  AtcoCode              NAME NOT NULL,
  IataCode              NAME NOT NULL,
  Name                  NAME NOT NULL,
  NameLang              NAME,
  CreationDateTime      NAME,
  ModificationDateTime  NAME,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (AtcoCode)
);

CREATE INDEX AirReferences_ncc ON naptan.AirReferences(IataCode);
CREATE INDEX AirReferences_n ON naptan.AirReferences(Name);

-- ================================================================================
-- AlternativeDescriptors
-- ================================================================================
DROP TABLE IF EXISTS naptan.AlternativeDescriptors CASCADE;

CREATE TABLE naptan.AlternativeDescriptors (
  AtcoCode              NAME NOT NULL,
  CommonName            NAME NOT NULL,
  CommonNameLang        NAME,
  ShortName             NAME,
  ShortCommonNameLang   NAME,
  Landmark              NAME,
  LandmarkLang          NAME,
  Street                NAME,
  StreetLang            NAME,
  Crossing              NAME,
  CrossingLang          NAME,
  Indicator             NAME,
  IndicatorLang         NAME,
  CreationDateTime      NAME,
  ModificationDateTime  NAME,
  RevisionNumber        NAME,
  Modification          NAME,
  PRIMARY KEY (AtcoCode,CommonName)
);

CREATE INDEX AlternativeDescriptors_ac ON naptan.AlternativeDescriptors(AtcoCode);
CREATE INDEX AlternativeDescriptors_cn ON naptan.AlternativeDescriptors(CommonName);
CREATE INDEX AlternativeDescriptors_sn ON naptan.AlternativeDescriptors(ShortName);

-- ================================================================================
-- AreaHierarchy
-- ================================================================================
DROP TABLE IF EXISTS naptan.AreaHierarchy CASCADE;

CREATE TABLE naptan.AreaHierarchy (
  ParentStopAreaCode    NAME NOT NULL,
  ChildStopAreaCode     NAME NOT NULL,
  CreationDateTime      NAME,
  ModificationDateTime  NAME,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (ParentStopAreaCode,ChildStopAreaCode)
);

CREATE INDEX AreaHierarchy_p ON naptan.AreaHierarchy(ParentStopAreaCode);
CREATE INDEX AreaHierarchy_c ON naptan.AreaHierarchy(ChildStopAreaCode);

-- ================================================================================
-- CoachReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.CoachReferences CASCADE;

CREATE TABLE naptan.CoachReferences (
  AtcoCode              NAME NOT NULL,
  OperatorRef           NAME NOT NULL,
  NationalCoachCode     NAME NOT NULL,
  Name                  NAME NOT NULL,
  NameLang              NAME NOT NULL,
  LongName              NAME NOT NULL,
  LongNameLang          NAME NOT NULL,
  GridType              NAME NOT NULL,
  Easting               INTEGER NOT NULL,
  Northing              INTEGER NOT NULL,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (AtcoCode)
);

CREATE INDEX CoachReferences_ncc ON naptan.CoachReferences(NationalCoachCode);
CREATE INDEX CoachReferences_n ON naptan.CoachReferences(Name);
CREATE INDEX CoachReferences_ln ON naptan.CoachReferences(LongName);

-- geometry
SELECT addgeometrycolumn( '', 'naptan', 'coachreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX CoachReferences_geom ON naptan.CoachReferences USING GIST (geom);

-- ================================================================================
-- FerryReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.FerryReferences CASCADE;

CREATE TABLE naptan.FerryReferences (
  AtcoCode              NAME NOT NULL,
  FerryCode             NAME NOT NULL,
  Name                  NAME NOT NULL,
  NameLang              NAME NOT NULL,
  GridType              NAME NOT NULL,
  Easting               INTEGER NOT NULL,
  Northing              INTEGER NOT NULL,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (AtcoCode)
);

CREATE INDEX FerryReferences_c ON naptan.FerryReferences(FerryCode);
CREATE INDEX FerryReferences_n ON naptan.FerryReferences(Name);

-- geometry
SELECT addgeometrycolumn( '', 'naptan', 'ferryreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX FerryReferences_geom ON naptan.FerryReferences USING GIST (geom);

-- ================================================================================
-- railreference
-- ================================================================================
DROP TABLE IF EXISTS naptan.RailReferences CASCADE;

CREATE TABLE naptan.RailReferences (
  AtcoCode              NAME NOT NULL,
  TiplocCode            NAME,
  CrsCode               CHAR(3),
  StationName           NAME,
  StationNameLang       NAME,
  GridType              NAME,
  Easting               INTEGER NOT NULL,
  Northing              INTEGER NOT NULL,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (AtcoCode)
);

CREATE INDEX RailReferences_t ON naptan.RailReferences(TiplocCode);
CREATE INDEX RailReferences_c ON naptan.RailReferences(CrsCode);
CREATE INDEX RailReferences_n ON naptan.RailReferences(lower(StationName));

-- geometry
SELECT addgeometrycolumn( '', 'naptan', 'railreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX RailReferences_geom ON naptan.RailReferences USING GIST (geom);

-- ================================================================================
-- stops
-- ================================================================================
DROP TABLE IF EXISTS naptan.stops CASCADE;

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

-- ================================================================================
-- stopplusbuszones - links naptan.stops with nptg.plusbus to allow us to filter
-- stops within a specific plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS naptan.stopplusbuszones CASCADE;

-- Note no references here as we may not have the entries present
CREATE TABLE naptan.stopplusbuszones (
  atco  NAME NOT NULL,
  zone  NAME NOT NULL,
  PRIMARY KEY (atco, zone)
);

CREATE INDEX stopplusbuszones_atco ON naptan.stopplusbuszones(atco);
CREATE INDEX stopplusbuszones_zone ON naptan.stopplusbuszones(zone);

-- ================================================================================
-- plusbusstops is a view of stops that only exist within a plusbus zone.
-- As this gets it's geometry from stops it can be used as a point feature
-- ================================================================================
CREATE VIEW naptan.plusbusstops
  AS SELECT z.zone AS plusbuszone, s.*
    FROM naptan.stops s
    INNER JOIN naptan.stopplusbuszones z ON s.atco = z.atco;
