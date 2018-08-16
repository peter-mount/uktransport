-- ================================================================================
-- nptg schema
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS nptg;

-- ================================================================================
-- PlusbusMapping consists of a polygon for each Plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS nptg.PlusbusMapping CASCADE;

CREATE TABLE nptg.PlusbusMapping (
  PlusbusZoneCode       NAME NOT NULL,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (PlusbusZoneCode)
);

-- geometry
SELECT addgeometrycolumn( '', 'nptg', 'plusbusmapping', 'geom', 27700, 'POLYGON', 2, true);
CREATE INDEX plusbusmapping_geom ON nptg.PlusbusMapping USING GIST (geom);

-- ================================================================================
-- PlusbusZones Holds the metadata for a Plusbus zone
-- ================================================================================
DROP TABLE IF EXISTS nptg.PlusbusZones CASCADE;

CREATE TABLE nptg.PlusbusZones (
  PlusbusZoneCode       NAME NOT NULL,
  Name                  NAME,
  NameLang              NAME,
  Country               NAME,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (PlusbusZoneCode)
);

-- ================================================================================
-- railplusbus is a view of naptan.rail that is contained within a plusbus zone
-- Specifically this means rail stations within each zone.
-- This view inherits the geometry from naptan.rail so can be used as a Point feature
-- ================================================================================
CREATE VIEW nptg.railplusbus
  AS SELECT
      z.PlusbusZoneCode,
      r.*
    FROM naptan.rail r
      INNER JOIN nptg.PlusbusMapping z ON ST_Contains( z.geom, r.geom );

-- ================================================================================
-- AdjacentLocality
-- ================================================================================
DROP TABLE IF EXISTS nptg.AdjacentLocality CASCADE;

CREATE TABLE nptg.AdjacentLocality (
  NptgLocalityCode          NAME NOT NULL,
  AdjacentNptgLocalityCode  NAME NOT NULL,
  CreationDateTime          TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber            INTEGER,
  Modification              NAME,
  PRIMARY KEY (NptgLocalityCode,AdjacentNptgLocalityCode)
);

CREATE INDEX AdjacentLocality_n ON nptg.AdjacentLocality( NptgLocalityCode );
CREATE INDEX AdjacentLocality_a ON nptg.AdjacentLocality( AdjacentNptgLocalityCode );

-- ================================================================================
-- AdminAreas
-- ================================================================================
DROP TABLE IF EXISTS nptg.AdminAreas CASCADE;

CREATE TABLE nptg.AdminAreas (
  AdministrativeAreaCode      NAME NOT NULL,
  AtcoAreaCode                NAME NOT NULL,
  AreaName                    NAME NOT NULL,
  AreaNameLang                NAME,
  ShortName                   NAME NOT NULL,
  ShortNameLang               NAME,
  Country                     NAME NOT NULL,
  RegionCode                  NAME,
  MaximumLengthForShortNames  NAME,
  National                    NAME,
  ContactEmail                NAME,
  ContactTelephone            NAME,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (AdministrativeAreaCode)
);

CREATE INDEX AdminAreas_a ON nptg.AdminAreas( AtcoAreaCode );
CREATE INDEX AdminAreas_n ON nptg.AdminAreas( AreaName );

-- ================================================================================
-- Districts
-- ================================================================================
DROP TABLE IF EXISTS nptg.Districts CASCADE;

CREATE TABLE nptg.Districts (
  DistrictCode            NAME NOT NULL,
  DistrictName            NAME NOT NULL,
  DistrictLang            NAME,
  AdministrativeAreaCode NAME NOT NULL,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (DistrictCode)
);

CREATE INDEX Districts_n ON nptg.Districts( DistrictName );
CREATE INDEX Districts_a ON nptg.Districts( AdministrativeAreaCode );

-- ================================================================================
-- Localities
-- ================================================================================
DROP TABLE IF EXISTS nptg.Localities CASCADE;

CREATE TABLE nptg.Localities (
  NptgLocalityCode        NAME NOT NULL,
  LocalityName            NAME NOT NULL,
  LocalityNameLang        NAME NOT NULL,
  ShortName               NAME NOT NULL,
  ShortNameLang           NAME NOT NULL,
  QualifierName           NAME NOT NULL,
  QualifierNameLang       NAME NOT NULL,
  QualifierLocalityRef    NAME NOT NULL,
  QualifierDistrictRef    NAME NOT NULL,
  AdministrativeAreaCode  NAME NOT NULL,
  NptgDistrictCode        NAME NOT NULL,
  SourceLocalityType      NAME NOT NULL,
  CreationDateTime        TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime    TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber          INTEGER,
  Modification            NAME,
  PRIMARY KEY (NptgLocalityCode)
);

-- geometry
SELECT addgeometrycolumn( '', 'nptg', 'localities', 'geom', 27700, 'POINT', 2, true);
CREATE INDEX Localities_geom ON nptg.Localities USING GIST (geom);

CREATE INDEX Localities_ln ON nptg.Localities( LocalityName );
CREATE INDEX Localities_sn ON nptg.Localities( ShortName );
CREATE INDEX Localities_qn ON nptg.Localities( QualifierName );
CREATE INDEX Localities_lr ON nptg.Localities( QualifierLocalityRef );
CREATE INDEX Localities_dr ON nptg.Localities( QualifierDistrictRef );
CREATE INDEX Localities_ac ON nptg.Localities( AdministrativeAreaCode );
CREATE INDEX Localities_dc ON nptg.Localities( NptgDistrictCode );

-- ================================================================================
-- LocalityAlternativeNames
-- ================================================================================
DROP TABLE IF EXISTS nptg.LocalityAlternativeNames CASCADE;

CREATE TABLE nptg.LocalityAlternativeNames (
  NptgLocalityCode        NAME NOT NULL,
  OldNptgLocalityCode     NAME NOT NULL,
  LocalityName            NAME NOT NULL,
  LocalityNameLang        NAME NOT NULL,
  ShortName               NAME,
  ShortNameLang           NAME NOT NULL,
  QualifierName           NAME,
  QualifierNameLang       NAME NOT NULL,
  QualifierLocalityRef    NAME,
  QualifierDistrictRef    NAME,
  CreationDateTime        TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime    TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber          INTEGER,
  Modification            NAME,
  PRIMARY KEY (NptgLocalityCode,OldNptgLocalityCode)
);

CREATE INDEX LocalityAlternativeNames_nc ON nptg.LocalityAlternativeNames( NptgLocalityCode );
CREATE INDEX LocalityAlternativeNames_oc ON nptg.LocalityAlternativeNames( OldNptgLocalityCode );
CREATE INDEX LocalityAlternativeNames_ln ON nptg.LocalityAlternativeNames( LocalityName );
CREATE INDEX LocalityAlternativeNames_sn ON nptg.LocalityAlternativeNames( ShortName );

-- ================================================================================
-- LocalityHierarchy
-- ================================================================================
DROP TABLE IF EXISTS nptg.LocalityHierarchy CASCADE;

CREATE TABLE nptg.LocalityHierarchy (
  ParentNptgLocalityCode  NAME NOT NULL,
  ChildNptgLocalityCode   NAME NOT NULL,
  CreationDateTime        TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime    TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber          INTEGER,
  Modification            NAME,
  PRIMARY KEY (ParentNptgLocalityCode,ChildNptgLocalityCode)
);

CREATE INDEX LocalityHierarchy_p ON nptg.LocalityHierarchy( ParentNptgLocalityCode );
CREATE INDEX LocalityHierarchy_c ON nptg.LocalityHierarchy( ChildNptgLocalityCode );

-- ================================================================================
-- Regions
-- ================================================================================
DROP TABLE IF EXISTS nptg.Regions CASCADE;

CREATE TABLE nptg.Regions (
  RegionCode            NAME NOT NULL,
  RegionName            NAME NOT NULL,
  RegionNameLang        NAME,
  CreationDateTime      TIMESTAMP WITHOUT TIME ZONE,
  ModificationDateTime  TIMESTAMP WITHOUT TIME ZONE,
  RevisionNumber        INTEGER,
  Modification          NAME,
  PRIMARY KEY (RegionCode)
);

CREATE INDEX Regions_n ON nptg.Regions( RegionName );
