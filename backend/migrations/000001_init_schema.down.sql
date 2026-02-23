-- =============================================================================
-- Migration: 000001_init_schema (DOWN)
-- Description: Completely reverse the init_schema migration.
--              Tables, types, extensions, and schemas are dropped in reverse
--              dependency order to avoid FK violations.
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1. Drop AI_ENGINE tables (reverse creation order)
-- ---------------------------------------------------------------------------
DROP TABLE IF EXISTS ai_engine.t_chat_messages;
DROP TABLE IF EXISTS ai_engine.t_chat_sessions;
DROP TABLE IF EXISTS ai_engine.t_communications;
DROP TABLE IF EXISTS ai_engine.t_hiring_forecasts;
DROP TABLE IF EXISTS ai_engine.t_processing_tasks;
DROP TABLE IF EXISTS ai_engine.t_resume_embeddings;
DROP TABLE IF EXISTS ai_engine.t_score_factors;
DROP TABLE IF EXISTS ai_engine.t_candidate_scores;

-- ---------------------------------------------------------------------------
-- 2. Drop HIRING tables (reverse creation order)
-- ---------------------------------------------------------------------------
DROP TABLE IF EXISTS hiring.t_activity_logs;
DROP TABLE IF EXISTS hiring.t_action_types;
DROP TABLE IF EXISTS hiring.t_status_history;
DROP TABLE IF EXISTS hiring.t_candidate_profiles;
DROP TABLE IF EXISTS hiring.t_candidates;
DROP TABLE IF EXISTS hiring.t_jobs;

-- ---------------------------------------------------------------------------
-- 3. Drop AUTH tables (reverse creation order)
-- ---------------------------------------------------------------------------
DROP TABLE IF EXISTS auth.t_users;
DROP TABLE IF EXISTS auth.t_teams;

-- ---------------------------------------------------------------------------
-- 4. Drop ENUM types
-- ---------------------------------------------------------------------------
DROP TYPE IF EXISTS message_role;
DROP TYPE IF EXISTS chat_type;
DROP TYPE IF EXISTS actor_type;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS email_type;
DROP TYPE IF EXISTS factor_type;
DROP TYPE IF EXISTS work_format_type;
DROP TYPE IF EXISTS user_role;

-- ---------------------------------------------------------------------------
-- 5. Drop extensions
-- ---------------------------------------------------------------------------
DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS "uuid-ossp";

-- ---------------------------------------------------------------------------
-- 6. Drop schemas
-- ---------------------------------------------------------------------------
DROP SCHEMA IF EXISTS ai_engine;
DROP SCHEMA IF EXISTS hiring;
DROP SCHEMA IF EXISTS auth;

COMMIT;
