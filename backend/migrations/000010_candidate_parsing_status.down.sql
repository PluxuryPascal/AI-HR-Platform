-- =============================================================================
-- Migration: 000010_candidate_parsing_status (DOWN)
-- Description: Remove parsing_status column and ENUM
-- =============================================================================

BEGIN;

ALTER TABLE hiring.t_candidates
    DROP COLUMN IF EXISTS parsing_status;

DROP TYPE IF EXISTS candidate_parsing_status;

COMMIT;
