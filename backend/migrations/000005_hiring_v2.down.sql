-- =============================================================================
-- Migration: 000005_hiring_v2 (DOWN)
-- Description: Reverse all changes from the UP migration, restoring the
--              original schema layout.
-- =============================================================================

BEGIN;

-- ===========================================================================
-- 9. Drop t_candidate_stage_history → recreate t_status_history
-- ===========================================================================
DROP TABLE IF EXISTS hiring.t_candidate_stage_history;

CREATE TABLE hiring.t_status_history (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id UUID NOT NULL REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    old_status   VARCHAR,
    new_status   VARCHAR,
    changed_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ===========================================================================
-- 8. Drop t_candidate_stages
-- ===========================================================================
DROP TABLE IF EXISTS hiring.t_candidate_stages;

-- ===========================================================================
-- 7. Drop t_pipeline_stages
-- ===========================================================================
DROP TABLE IF EXISTS hiring.t_pipeline_stages;

-- ===========================================================================
-- 6. Revert t_candidate_profiles — drop ai_parsed_at
-- ===========================================================================
ALTER TABLE hiring.t_candidate_profiles
    DROP COLUMN IF EXISTS ai_parsed_at;

-- ===========================================================================
-- 5. Restore t_candidates columns
-- ===========================================================================
ALTER TABLE hiring.t_candidates
    ADD COLUMN status          VARCHAR,
    ADD COLUMN kanban_position FLOAT;

-- ===========================================================================
-- 4. Revert t_jobs
-- ===========================================================================

-- Convert status back to VARCHAR
ALTER TABLE hiring.t_jobs
    ALTER COLUMN status DROP NOT NULL,
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE VARCHAR USING status::TEXT;

ALTER TABLE hiring.t_jobs
    DROP COLUMN IF EXISTS salary_min,
    DROP COLUMN IF EXISTS salary_max,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS updated_at;

-- ===========================================================================
-- 3. Restore cross-schema foreign keys
-- ===========================================================================

-- 3a. hiring.t_jobs → auth.t_teams
ALTER TABLE hiring.t_jobs
    ADD CONSTRAINT t_jobs_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES auth.t_teams (id) ON DELETE CASCADE;

-- 3b. hiring.t_job_access → auth.t_users
ALTER TABLE hiring.t_job_access
    ADD CONSTRAINT t_job_access_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES auth.t_users (id) ON DELETE CASCADE;

-- 3c. auth.t_invite_job_access → hiring.t_jobs
ALTER TABLE auth.t_invite_job_access
    ADD CONSTRAINT t_invite_job_access_job_id_fkey
    FOREIGN KEY (job_id) REFERENCES hiring.t_jobs (id) ON DELETE CASCADE;

-- 3d. ai_engine.t_candidate_scores → hiring.t_candidates
ALTER TABLE ai_engine.t_candidate_scores
    ADD CONSTRAINT t_candidate_scores_candidate_id_fkey
    FOREIGN KEY (candidate_id) REFERENCES hiring.t_candidates (id) ON DELETE CASCADE;

-- 3e. ai_engine.t_resume_embeddings → auth.t_teams, hiring.t_candidates
ALTER TABLE ai_engine.t_resume_embeddings
    ADD CONSTRAINT t_resume_embeddings_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    ADD CONSTRAINT t_resume_embeddings_candidate_id_fkey
    FOREIGN KEY (candidate_id) REFERENCES hiring.t_candidates (id) ON DELETE CASCADE;

-- 3f. ai_engine.t_hiring_forecasts → auth.t_teams
ALTER TABLE ai_engine.t_hiring_forecasts
    ADD CONSTRAINT t_hiring_forecasts_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES auth.t_teams (id) ON DELETE CASCADE;

-- 3g. ai_engine.t_communications → hiring.t_candidates, auth.t_users
ALTER TABLE ai_engine.t_communications
    ADD CONSTRAINT t_communications_candidate_id_fkey
    FOREIGN KEY (candidate_id) REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    ADD CONSTRAINT t_communications_generated_by_user_id_fkey
    FOREIGN KEY (generated_by_user_id) REFERENCES auth.t_users (id) ON DELETE CASCADE;

-- 3h. ai_engine.t_chat_sessions → auth.t_teams, auth.t_users, hiring.t_candidates
ALTER TABLE ai_engine.t_chat_sessions
    ADD CONSTRAINT t_chat_sessions_team_id_fkey
    FOREIGN KEY (team_id) REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    ADD CONSTRAINT t_chat_sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES auth.t_users (id) ON DELETE CASCADE,
    ADD CONSTRAINT t_chat_sessions_target_candidate_id_fkey
    FOREIGN KEY (target_candidate_id) REFERENCES hiring.t_candidates (id) ON DELETE SET NULL;

-- ===========================================================================
-- 2. Restore hiring.t_action_types & t_activity_logs from logging
-- ===========================================================================
CREATE TABLE hiring.t_action_types (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE hiring.t_activity_logs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id    UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    actor_type actor_type NOT NULL,
    actor_id   UUID,
    action_id  INT NOT NULL REFERENCES hiring.t_action_types (id) ON DELETE RESTRICT,
    target_id  UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Migrate data back from logging schema
INSERT INTO hiring.t_action_types (id, code, description)
SELECT id, code, description
FROM logging.t_action_types
WHERE service = 'hiring';

SELECT setval(
    pg_get_serial_sequence('hiring.t_action_types', 'id'),
    COALESCE((SELECT MAX(id) FROM hiring.t_action_types), 0)
);

INSERT INTO hiring.t_activity_logs
    (id, team_id, actor_type, actor_id, action_id, target_id, created_at)
SELECT id, team_id, actor_type, actor_id, action_id, target_id, created_at
FROM logging.t_activity_logs
WHERE service = 'hiring';

-- Drop logging tables & schema
DROP TABLE logging.t_activity_logs;
DROP TABLE logging.t_action_types;
DROP SCHEMA IF EXISTS logging;

-- ===========================================================================
-- 1. Drop new types
-- ===========================================================================
DROP TYPE IF EXISTS job_status;

COMMIT;
