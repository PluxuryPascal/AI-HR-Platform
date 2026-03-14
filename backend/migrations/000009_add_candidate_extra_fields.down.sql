-- =============================================================================
-- Migration: 000009_add_candidate_extra_fields (DOWN)
-- Description: Remove location and skills fields from t_candidates
-- =============================================================================

BEGIN;

ALTER TABLE hiring.t_candidates
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS skills;

COMMIT;
