-- =============================================================================
-- Migration: 000011_add_candidate_profile_missing_fields (DOWN)
-- Description: Remove missing_fields column from t_candidate_profiles
-- =============================================================================

BEGIN;

ALTER TABLE hiring.t_candidate_profiles
    DROP COLUMN IF EXISTS missing_fields;

COMMIT;
