ALTER TABLE ai_engine.t_communications
    ADD COLUMN subject TEXT,
    ADD COLUMN body    TEXT;

ALTER TABLE ai_engine.t_communications
    ALTER COLUMN content DROP NOT NULL;

CREATE INDEX idx_communications_candidate_id
    ON ai_engine.t_communications (candidate_id);
