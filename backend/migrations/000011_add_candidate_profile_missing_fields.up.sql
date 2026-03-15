-- =============================================================================
-- Migration: 000011_add_candidate_profile_missing_fields (UP)
-- Description: Add missing_fields column to t_candidate_profiles
-- =============================================================================

BEGIN;

ALTER TABLE hiring.t_candidate_profiles
    ADD COLUMN missing_fields TEXT[] NOT NULL DEFAULT '{}';

COMMIT;
