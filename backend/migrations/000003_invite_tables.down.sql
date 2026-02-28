-- =============================================================================
-- Migration: 000003_invite_tables (DOWN)
-- Description: Rollback — drop tables in reverse dependency order.
-- =============================================================================

BEGIN;

DROP TABLE IF EXISTS auth.t_invite_job_access;
DROP TABLE IF EXISTS hiring.t_job_access;
DROP TABLE IF EXISTS auth.t_invites;

COMMIT;
