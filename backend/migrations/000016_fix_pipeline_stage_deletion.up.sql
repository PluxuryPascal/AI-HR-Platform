BEGIN;

-- 1. Update t_candidate_stage_history (from_stage_id -> SET NULL)
ALTER TABLE hiring.t_candidate_stage_history
    DROP CONSTRAINT IF EXISTS t_candidate_stage_history_from_stage_id_fkey;

ALTER TABLE hiring.t_candidate_stage_history
    ADD CONSTRAINT t_candidate_stage_history_from_stage_id_fkey
    FOREIGN KEY (from_stage_id) REFERENCES hiring.t_pipeline_stages (id) ON DELETE SET NULL;

-- 2. Update t_candidate_stage_history (to_stage_id -> CASCADE)
-- If a stage is deleted, history records pointing TO it as the destination are removed.
ALTER TABLE hiring.t_candidate_stage_history
    DROP CONSTRAINT IF EXISTS t_candidate_stage_history_to_stage_id_fkey;

ALTER TABLE hiring.t_candidate_stage_history
    ADD CONSTRAINT t_candidate_stage_history_to_stage_id_fkey
    FOREIGN KEY (to_stage_id) REFERENCES hiring.t_pipeline_stages (id) ON DELETE CASCADE;

-- 3. Update t_candidate_stages (stage_id -> CASCADE)
-- If a stage is deleted, candidates are moved in application logic, 
-- but this constraint ensures SQL-level consistency.
ALTER TABLE hiring.t_candidate_stages
    DROP CONSTRAINT IF EXISTS t_candidate_stages_stage_id_fkey;

ALTER TABLE hiring.t_candidate_stages
    ADD CONSTRAINT t_candidate_stages_stage_id_fkey
    FOREIGN KEY (stage_id) REFERENCES hiring.t_pipeline_stages (id) ON DELETE CASCADE;

COMMIT;
