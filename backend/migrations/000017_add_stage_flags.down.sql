BEGIN;

ALTER TABLE hiring.t_pipeline_stages
    DROP COLUMN IF EXISTS is_rejection,
    DROP COLUMN IF EXISTS is_interview;

COMMIT;
