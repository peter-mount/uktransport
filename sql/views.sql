-- ================================================================================
-- views contains SQL views that act across multiple schemas
-- ================================================================================

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS views;

-- ================================================================================
-- railplusbus is a view of naptan.RailReferences that is contained within
-- a plusbus zone.
--
-- Specifically this means rail stations within each zone.
-- This view inherits the geometry from naptan.RailReferences so it can be
-- used as a Point feature
-- ================================================================================
CREATE OR REPLACE VIEW views.railplusbus
  AS SELECT
      z.PlusbusZoneCode,
      r.*
    FROM naptan.RailReferences r
      INNER JOIN nptg.PlusbusMapping z ON ST_Contains( z.geom, r.geom );
