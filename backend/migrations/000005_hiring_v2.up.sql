-- =============================================================================
-- Migration: 000005_hiring_v2 (UP)
-- Description: Restructure the hiring schema for the Hiring microservice:
--   - Introduce logging schema (centralised activity logs)
--   - Add job_status ENUM, pipeline stages, candidate stage tracking
--   - Add salary fields, created_by to t_jobs
--   - Move kanban_position/status out of t_candidates into t_candidate_stages
--   - Replace t_status_history with t_candidate_stage_history
--   - Remove all cross-schema foreign keys (soft UUID refs)
-- =============================================================================

BEGIN;

-- ===========================================================================
-- 1. New schemas & types
-- ===========================================================================
CREATE SCHEMA IF NOT EXISTS logging;

CREATE TYPE job_status AS ENUM (
    'status_draft',
    'status_published',
    'status_closed',
    'status_archived'
);

-- ===========================================================================
-- 2. Centralised logging (moved from hiring.t_action_types / t_activity_logs)
-- ===========================================================================

-- 2a. New logging tables ────────────────────────────────────────────────────
CREATE TABLE logging.t_action_types (
    id          SERIAL PRIMARY KEY,
    service     VARCHAR NOT NULL,          -- 'auth', 'hiring', 'ai_engine'
    code        VARCHAR NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE logging.t_activity_logs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id    UUID NOT NULL,              -- soft ref (no FK)
    service    VARCHAR NOT NULL,           -- source service
    actor_type actor_type NOT NULL,
    actor_id   UUID,
    action_id  INT NOT NULL REFERENCES logging.t_action_types (id) ON DELETE RESTRICT,
    target_id  UUID,
    payload    JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_team_created
    ON logging.t_activity_logs (team_id, created_at DESC);
CREATE INDEX idx_activity_target
    ON logging.t_activity_logs (target_id);

-- 2b. Migrate existing data ────────────────────────────────────────────────

-- Action types: copy rows, adding default service = 'hiring'
INSERT INTO logging.t_action_types (id, service, code, description)
SELECT id, 'hiring', code, description
FROM hiring.t_action_types;

-- Sync the SERIAL sequence so future inserts get correct ids
SELECT setval(
    pg_get_serial_sequence('logging.t_action_types', 'id'),
    COALESCE((SELECT MAX(id) FROM logging.t_action_types), 1),
    (SELECT COUNT(*) > 0 FROM logging.t_action_types)
);

-- Activity logs: copy rows, adding default service = 'hiring'
INSERT INTO logging.t_activity_logs
    (id, team_id, service, actor_type, actor_id, action_id, target_id, created_at)
SELECT id, team_id, 'hiring', actor_type, actor_id, action_id, target_id, created_at
FROM hiring.t_activity_logs;

-- 2c. Drop old hiring logging tables ──────────────────────────────────────
DROP TABLE hiring.t_activity_logs;
DROP TABLE hiring.t_action_types;

-- ===========================================================================
-- 3. Remove cross-schema foreign keys
-- ===========================================================================

-- 3a. hiring.t_jobs → auth.t_teams
ALTER TABLE hiring.t_jobs
    DROP CONSTRAINT t_jobs_team_id_fkey;

-- 3b. hiring.t_job_access → auth.t_users
ALTER TABLE hiring.t_job_access
    DROP CONSTRAINT t_job_access_user_id_fkey;

-- 3c. auth.t_invite_job_access → hiring.t_jobs
ALTER TABLE auth.t_invite_job_access
    DROP CONSTRAINT t_invite_job_access_job_id_fkey;

-- 3d. ai_engine.t_candidate_scores → hiring.t_candidates
ALTER TABLE ai_engine.t_candidate_scores
    DROP CONSTRAINT t_candidate_scores_candidate_id_fkey;

-- 3e. ai_engine.t_resume_embeddings → auth.t_teams, hiring.t_candidates
ALTER TABLE ai_engine.t_resume_embeddings
    DROP CONSTRAINT t_resume_embeddings_team_id_fkey,
    DROP CONSTRAINT t_resume_embeddings_candidate_id_fkey;

-- 3f. ai_engine.t_hiring_forecasts → auth.t_teams
ALTER TABLE ai_engine.t_hiring_forecasts
    DROP CONSTRAINT t_hiring_forecasts_team_id_fkey;

-- 3g. ai_engine.t_communications → hiring.t_candidates, auth.t_users
ALTER TABLE ai_engine.t_communications
    DROP CONSTRAINT t_communications_candidate_id_fkey,
    DROP CONSTRAINT t_communications_generated_by_user_id_fkey;

-- 3h. ai_engine.t_chat_sessions → auth.t_teams, auth.t_users, hiring.t_candidates
ALTER TABLE ai_engine.t_chat_sessions
    DROP CONSTRAINT t_chat_sessions_team_id_fkey,
    DROP CONSTRAINT t_chat_sessions_user_id_fkey,
    DROP CONSTRAINT t_chat_sessions_target_candidate_id_fkey;

-- 3i. hiring.t_candidate_profiles → hiring.t_candidates (keep WITHIN hiring)
--     This one stays — both tables belong to the same service.

-- ===========================================================================
-- 4. Alter hiring.t_jobs
-- ===========================================================================
ALTER TABLE hiring.t_jobs
    ADD COLUMN salary_min  INT,
    ADD COLUMN salary_max  INT,
    ADD COLUMN currency    VARCHAR(3) DEFAULT 'RUB',
    ADD COLUMN created_by  UUID,                        -- soft ref to auth.t_users
    ADD COLUMN updated_at  TIMESTAMP NOT NULL DEFAULT NOW();

-- Convert status VARCHAR → job_status ENUM
-- Existing rows with NULL or unrecognised values default to 'status_draft'
ALTER TABLE hiring.t_jobs
    ALTER COLUMN status TYPE job_status
        USING CASE
            WHEN status = 'status_published' THEN 'status_published'::job_status
            WHEN status = 'status_closed'    THEN 'status_closed'::job_status
            WHEN status = 'status_archived'  THEN 'status_archived'::job_status
            ELSE 'status_draft'::job_status
        END;

ALTER TABLE hiring.t_jobs
    ALTER COLUMN status SET DEFAULT 'status_draft',
    ALTER COLUMN status SET NOT NULL;

-- ===========================================================================
-- 5. Alter hiring.t_candidates — remove inline status & kanban_position
-- ===========================================================================
ALTER TABLE hiring.t_candidates
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS kanban_position;

-- ===========================================================================
-- 6. Alter hiring.t_candidate_profiles — add AI timestamp
-- ===========================================================================
ALTER TABLE hiring.t_candidate_profiles
    ADD COLUMN ai_parsed_at TIMESTAMP;

-- ===========================================================================
-- 7. Create hiring.t_pipeline_stages
-- ===========================================================================
CREATE TABLE hiring.t_pipeline_stages (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id      UUID,                                     -- NULL = team template
    team_id     UUID NOT NULL,                            -- soft ref to auth.t_teams
    code        VARCHAR NOT NULL,
    title       VARCHAR NOT NULL,
    "position"  INT NOT NULL,
    is_terminal BOOLEAN NOT NULL DEFAULT FALSE,
    color       VARCHAR
);

-- Expression-based unique: one code per pipeline (per job or per team-template)
CREATE UNIQUE INDEX uq_stage_code
    ON hiring.t_pipeline_stages (
        COALESCE(job_id, '00000000-0000-0000-0000-000000000000'::UUID),
        team_id,
        code
    );

-- job_id FK is within hiring schema — allowed
ALTER TABLE hiring.t_pipeline_stages
    ADD CONSTRAINT fk_pipeline_stages_job
    FOREIGN KEY (job_id) REFERENCES hiring.t_jobs (id) ON DELETE CASCADE;

CREATE INDEX idx_pipeline_stages_job  ON hiring.t_pipeline_stages (job_id);
CREATE INDEX idx_pipeline_stages_team ON hiring.t_pipeline_stages (team_id);

-- ===========================================================================
-- 8. Create hiring.t_candidate_stages (1:1 — current position on kanban)
-- ===========================================================================
CREATE TABLE hiring.t_candidate_stages (
    candidate_id UUID PRIMARY KEY REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    stage_id     UUID NOT NULL    REFERENCES hiring.t_pipeline_stages (id),
    "position"   FLOAT NOT NULL DEFAULT 0,
    moved_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_candidate_stages_stage ON hiring.t_candidate_stages (stage_id);

-- ===========================================================================
-- 9. Replace t_status_history → t_candidate_stage_history
-- ===========================================================================
DROP TABLE hiring.t_status_history;

CREATE TABLE hiring.t_candidate_stage_history (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id  UUID NOT NULL REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    from_stage_id UUID REFERENCES hiring.t_pipeline_stages (id),
    to_stage_id   UUID NOT NULL REFERENCES hiring.t_pipeline_stages (id),
    changed_by    UUID NOT NULL,           -- soft ref to auth.t_users
    changed_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stage_history_candidate
    ON hiring.t_candidate_stage_history (candidate_id, changed_at DESC);

COMMIT;
