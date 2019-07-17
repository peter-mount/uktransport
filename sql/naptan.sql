-- ================================================================================
-- naptan
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS naptan;

-- ================================================================================
-- AirReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.AirReferences CASCADE;

CREATE TABLE naptan.AirReferences
(
    AtcoCode             NAME NOT NULL,
    IataCode             NAME NOT NULL,
    Name                 NAME NOT NULL,
    NameLang             NAME,
    CreationDateTime     NAME,
    ModificationDateTime NAME,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (AtcoCode)
);

CREATE INDEX AirReferences_ncc ON naptan.AirReferences (IataCode);
CREATE INDEX AirReferences_n ON naptan.AirReferences (Name);

-- ================================================================================
-- AlternativeDescriptors
-- ================================================================================
DROP TABLE IF EXISTS naptan.AlternativeDescriptors CASCADE;

CREATE TABLE naptan.AlternativeDescriptors
(
    AtcoCode             NAME NOT NULL,
    CommonName           NAME NOT NULL,
    CommonNameLang       NAME,
    ShortName            NAME,
    ShortCommonNameLang  NAME,
    Landmark             NAME,
    LandmarkLang         NAME,
    Street               NAME,
    StreetLang           NAME,
    Crossing             NAME,
    CrossingLang         NAME,
    Indicator            NAME,
    IndicatorLang        NAME,
    CreationDateTime     NAME,
    ModificationDateTime NAME,
    RevisionNumber       NAME,
    Modification         NAME,
    PRIMARY KEY (AtcoCode, CommonName)
);

CREATE INDEX AlternativeDescriptors_ac ON naptan.AlternativeDescriptors (AtcoCode);
CREATE INDEX AlternativeDescriptors_cn ON naptan.AlternativeDescriptors (CommonName);
CREATE INDEX AlternativeDescriptors_sn ON naptan.AlternativeDescriptors (ShortName);

-- ================================================================================
-- AreaHierarchy
-- ================================================================================
DROP TABLE IF EXISTS naptan.AreaHierarchy CASCADE;

CREATE TABLE naptan.AreaHierarchy
(
    ParentStopAreaCode   NAME NOT NULL,
    ChildStopAreaCode    NAME NOT NULL,
    CreationDateTime     NAME,
    ModificationDateTime NAME,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (ParentStopAreaCode, ChildStopAreaCode)
);

CREATE INDEX AreaHierarchy_p ON naptan.AreaHierarchy (ParentStopAreaCode);
CREATE INDEX AreaHierarchy_c ON naptan.AreaHierarchy (ChildStopAreaCode);

-- ================================================================================
-- CoachReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.CoachReferences CASCADE;

CREATE TABLE naptan.CoachReferences
(
    AtcoCode             NAME    NOT NULL,
    OperatorRef          NAME    NOT NULL,
    NationalCoachCode    NAME    NOT NULL,
    Name                 NAME    NOT NULL,
    NameLang             NAME    NOT NULL,
    LongName             NAME    NOT NULL,
    LongNameLang         NAME    NOT NULL,
    GridType             NAME    NOT NULL,
    Easting              INTEGER NOT NULL,
    Northing             INTEGER NOT NULL,
    CreationDateTime     TIMESTAMP WITHOUT TIME ZONE,
    ModificationDateTime TIMESTAMP WITHOUT TIME ZONE,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (AtcoCode)
);

CREATE INDEX CoachReferences_ncc ON naptan.CoachReferences (NationalCoachCode);
CREATE INDEX CoachReferences_n ON naptan.CoachReferences (Name);
CREATE INDEX CoachReferences_ln ON naptan.CoachReferences (LongName);

-- geometry
SELECT addgeometrycolumn('', 'naptan', 'coachreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX CoachReferences_geom ON naptan.CoachReferences USING GIST (geom);

-- ================================================================================
-- FerryReferences
-- ================================================================================
DROP TABLE IF EXISTS naptan.FerryReferences CASCADE;

CREATE TABLE naptan.FerryReferences
(
    AtcoCode             NAME    NOT NULL,
    FerryCode            NAME    NOT NULL,
    Name                 NAME    NOT NULL,
    NameLang             NAME    NOT NULL,
    GridType             NAME    NOT NULL,
    Easting              INTEGER NOT NULL,
    Northing             INTEGER NOT NULL,
    CreationDateTime     TIMESTAMP WITHOUT TIME ZONE,
    ModificationDateTime TIMESTAMP WITHOUT TIME ZONE,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (AtcoCode)
);

CREATE INDEX FerryReferences_c ON naptan.FerryReferences (FerryCode);
CREATE INDEX FerryReferences_n ON naptan.FerryReferences (Name);

-- geometry
SELECT addgeometrycolumn('', 'naptan', 'ferryreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX FerryReferences_geom ON naptan.FerryReferences USING GIST (geom);

-- ================================================================================
-- railreference
-- ================================================================================
DROP TABLE IF EXISTS naptan.RailReferences CASCADE;

CREATE TABLE naptan.RailReferences
(
    AtcoCode             NAME    NOT NULL,
    TiplocCode           NAME,
    CrsCode              CHAR(3),
    StationName          NAME,
    StationNameLang      NAME,
    GridType             NAME,
    Easting              INTEGER NOT NULL,
    Northing             INTEGER NOT NULL,
    CreationDateTime     TIMESTAMP WITHOUT TIME ZONE,
    ModificationDateTime TIMESTAMP WITHOUT TIME ZONE,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (AtcoCode)
);

CREATE INDEX RailReferences_t ON naptan.RailReferences (TiplocCode);
CREATE INDEX RailReferences_c ON naptan.RailReferences (CrsCode);
CREATE INDEX RailReferences_n ON naptan.RailReferences (lower(StationName));

-- geometry
SELECT addgeometrycolumn('', 'naptan', 'railreferences', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX RailReferences_geom ON naptan.RailReferences USING GIST (geom);

create or replace function naptan.railstationcrs(pcrs text)
    returns json
as
$$
SELECT row_to_json(t.*)
FROM naptan.railreferences t
WHERE crscode = pcrs
$$
    language sql;

create or replace function naptan.railstationtiploc(ptiploc text)
    returns json
as
$$
SELECT row_to_json(t.*)
FROM naptan.railreferences t
WHERE tiploccode = ptiploc
$$
    language sql;

create or replace function naptan.railstationatco(patco text)
    returns json
as
$$
SELECT row_to_json(t.*)
FROM naptan.railreferences t
WHERE atcocode = patco
$$
    language sql;

-- ================================================================================
-- stops
-- ================================================================================
DROP TABLE IF EXISTS naptan.Stops CASCADE;

CREATE TABLE naptan.Stops
(
    ATCOCode                NAME NOT NULL,
    NaptanCode              NAME,
    PlateCode               NAME,
    CleardownCode           NAME,
    CommonName              NAME,
    CommonNameLang          NAME,
    ShortCommonName         NAME,
    ShortCommonNameLang     NAME,
    Landmark                NAME,
    LandmarkLang            NAME,
    Street                  NAME,
    StreetLang              NAME,
    Crossing                NAME,
    CrossingLang            NAME,
    Indicator               NAME,
    IndicatorLang           NAME,
    Bearing                 NAME,
    NptgLocalityCode        NAME,
    LocalityName            NAME,
    ParentLocalityName      NAME,
    GrandParentLocalityName NAME,
    Town                    NAME,
    TownLang                NAME,
    Suburb                  NAME,
    SuburbLang              NAME,
    LocalityCentre          NAME,
    GridType                NAME,
    Easting                 INTEGER,
    Northing                INTEGER,
    Longitude               REAL,
    Latitude                REAL,
    StopType                NAME,
    BusStopType             NAME,
    TimingStatus            NAME,
    DefaultWaitTime         NAME,
    Notes                   NAME,
    NotesLang               NAME,
    AdministrativeAreaCode  NAME,
    CreationDateTime        TIMESTAMP WITHOUT TIME ZONE,
    ModificationDateTime    TIMESTAMP WITHOUT TIME ZONE,
    RevisionNumber          INTEGER,
    Modification            NAME,
    Status                  NAME,
    PRIMARY KEY (ATCOCode)
);

CREATE INDEX stops_nc ON naptan.stops (NaptanCode);
CREATE INDEX stops_cn ON naptan.stops (lower(CommonName));

-- geometry
SELECT addgeometrycolumn('', 'naptan', 'stops', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX stops_geom ON naptan.stops USING GIST (geom);

-- ================================================================================
-- StopPlusbusZones - links naptan.stops with nptg.plusbus to allow us to filter
-- stops within a specific plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS naptan.StopPlusbusZones CASCADE;

-- Note no references here as we may not have the entries present
CREATE TABLE naptan.StopPlusbusZones
(
    AtcoCode             NAME NOT NULL,
    PlusbusZoneCode      NAME NOT NULL,
    CreationDateTime     TIMESTAMP WITHOUT TIME ZONE,
    ModificationDateTime TIMESTAMP WITHOUT TIME ZONE,
    RevisionNumber       INTEGER,
    Modification         NAME,
    PRIMARY KEY (AtcoCode, PlusbusZoneCode)
);

CREATE INDEX stopplusbuszones_a ON naptan.stopplusbuszones (AtcoCode);
CREATE INDEX stopplusbuszones_p ON naptan.stopplusbuszones (PlusbusZoneCode);

-- ================================================================================
-- plusbusstops is a view of stops that only exist within a plusbus zone.
-- As this gets it's geometry from stops it can be used as a point feature
-- ================================================================================
CREATE VIEW naptan.plusbusstops
AS
SELECT z.PlusbusZoneCode, s.*
FROM naptan.Stops s
         INNER JOIN naptan.StopPlusbusZones z ON s.ATCOCode = z.AtcoCode;
