-- Migration: 000014_add_interview_guide (DOWN)
BEGIN;

ALTER TABLE hiring.t_candidate_profiles
    DROP COLUMN IF EXISTS interview_guide;

COMMIT;
