-- =============================================================================
-- Migration: 000009_add_candidate_extra_fields (UP)
-- Description: Add location and skills fields to t_candidates
-- =============================================================================

BEGIN;

ALTER TABLE hiring.t_candidates
    ADD COLUMN location VARCHAR,
    ADD COLUMN skills   TEXT[];

COMMIT;
