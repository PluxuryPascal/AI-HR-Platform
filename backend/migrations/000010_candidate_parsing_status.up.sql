-- =============================================================================
-- Migration: 000010_candidate_parsing_status (UP)
-- Description: Add parsing_status ENUM and column to t_candidates
-- =============================================================================

BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'candidate_parsing_status') THEN
        CREATE TYPE candidate_parsing_status AS ENUM (
            'pending',
            'processing',
            'needs_review',
            'completed',
            'failed'
        );
    END IF;
END
$$;

ALTER TABLE hiring.t_candidates
    ADD COLUMN IF NOT EXISTS parsing_status candidate_parsing_status NOT NULL DEFAULT 'pending';

COMMIT;
