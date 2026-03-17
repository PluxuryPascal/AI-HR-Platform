DROP INDEX IF EXISTS ai_engine.idx_communications_candidate_id;

ALTER TABLE ai_engine.t_communications
    ALTER COLUMN content SET NOT NULL;

ALTER TABLE ai_engine.t_communications
    DROP COLUMN IF EXISTS body,
    DROP COLUMN IF EXISTS subject;
