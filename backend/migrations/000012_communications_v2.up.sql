ALTER TABLE ai_engine.t_communications
    ADD COLUMN IF NOT EXISTS subject TEXT,
    ADD COLUMN IF NOT EXISTS body    TEXT;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'ai_engine' AND table_name = 't_communications'
          AND column_name = 'content' AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE ai_engine.t_communications ALTER COLUMN content DROP NOT NULL;
    END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_communications_candidate_id
    ON ai_engine.t_communications (candidate_id);
