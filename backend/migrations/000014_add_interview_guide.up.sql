-- Migration: 000014_add_interview_guide (UP)
BEGIN;

ALTER TABLE hiring.t_candidate_profiles
    ADD COLUMN IF NOT EXISTS interview_guide JSONB;

COMMIT;
