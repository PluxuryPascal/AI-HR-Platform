BEGIN;

ALTER TABLE hiring.t_candidate_stage_history
    DROP CONSTRAINT IF EXISTS t_candidate_stage_history_from_stage_id_fkey;
ALTER TABLE hiring.t_candidate_stage_history
    ADD CONSTRAINT t_candidate_stage_history_from_stage_id_fkey
    FOREIGN KEY (from_stage_id) REFERENCES hiring.t_pipeline_stages (id);

ALTER TABLE hiring.t_candidate_stage_history
    DROP CONSTRAINT IF EXISTS t_candidate_stage_history_to_stage_id_fkey;
ALTER TABLE hiring.t_candidate_stage_history
    ADD CONSTRAINT t_candidate_stage_history_to_stage_id_fkey
    FOREIGN KEY (to_stage_id) REFERENCES hiring.t_pipeline_stages (id);

ALTER TABLE hiring.t_candidate_stages
    DROP CONSTRAINT IF EXISTS t_candidate_stages_stage_id_fkey;
ALTER TABLE hiring.t_candidate_stages
    ADD CONSTRAINT t_candidate_stages_stage_id_fkey
    FOREIGN KEY (stage_id) REFERENCES hiring.t_pipeline_stages (id);

COMMIT;
