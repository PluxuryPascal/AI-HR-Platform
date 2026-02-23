-- =============================================================================
-- Migration: 000001_init_schema (UP)
-- Description: Initialise the full database schema for the Resume Screener
--              platform across three schemas: auth, hiring, ai_engine.
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1. Schemas
-- ---------------------------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS hiring;
CREATE SCHEMA IF NOT EXISTS ai_engine;

-- ---------------------------------------------------------------------------
-- 2. Extensions
-- ---------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector; -- pgvector for RAG embeddings

-- ---------------------------------------------------------------------------
-- 3. ENUM types (public schema so they are accessible from all schemas)
-- ---------------------------------------------------------------------------
CREATE TYPE user_role        AS ENUM ('admin', 'recruiter', 'hiring_manager');
CREATE TYPE work_format_type AS ENUM ('remote', 'office', 'hybrid');
CREATE TYPE factor_type      AS ENUM ('positive', 'negative', 'neutral');
CREATE TYPE email_type       AS ENUM ('rejection', 'interview_invite');
CREATE TYPE task_status      AS ENUM ('pending', 'extracting', 'analyzing', 'completed', 'failed');
CREATE TYPE actor_type       AS ENUM ('ai_agent', 'system', 'user');
CREATE TYPE chat_type        AS ENUM ('local_candidate', 'global_search');
CREATE TYPE message_role     AS ENUM ('user', 'assistant', 'system');

-- ===========================================================================
-- 4. AUTH schema tables
-- ===========================================================================

-- t_teams --------------------------------------------------------------------
CREATE TABLE auth.t_teams (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- t_users --------------------------------------------------------------------
CREATE TABLE auth.t_users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id       UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    email         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    role          user_role NOT NULL,
    locale        VARCHAR,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ===========================================================================
-- 5. HIRING schema tables
-- ===========================================================================

-- t_jobs ---------------------------------------------------------------------
CREATE TABLE hiring.t_jobs (
    id                     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id                UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    title                  VARCHAR NOT NULL,
    department             VARCHAR,
    work_format            work_format_type,
    description            TEXT,
    extracted_requirements JSONB,
    status                 VARCHAR,
    created_at             TIMESTAMP NOT NULL DEFAULT NOW()
);

-- t_candidates ---------------------------------------------------------------
CREATE TABLE hiring.t_candidates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id          UUID NOT NULL REFERENCES hiring.t_jobs (id) ON DELETE CASCADE,
    first_name      VARCHAR,
    last_name       VARCHAR,
    email           VARCHAR,
    resume_file_key VARCHAR,
    parsed_text     TEXT,
    status          VARCHAR,
    kanban_position FLOAT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP
);

-- t_candidate_profiles -------------------------------------------------------
CREATE TABLE hiring.t_candidate_profiles (
    candidate_id    UUID PRIMARY KEY REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    structured_data JSONB,
    updated_at      TIMESTAMP
);

-- t_status_history -----------------------------------------------------------
CREATE TABLE hiring.t_status_history (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id UUID NOT NULL REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    old_status   VARCHAR,
    new_status   VARCHAR,
    changed_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- t_action_types -------------------------------------------------------------
CREATE TABLE hiring.t_action_types (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR NOT NULL UNIQUE,
    description TEXT
);

-- t_activity_logs ------------------------------------------------------------
CREATE TABLE hiring.t_activity_logs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id    UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    actor_type actor_type NOT NULL,
    actor_id   UUID,
    action_id  INT NOT NULL REFERENCES hiring.t_action_types (id) ON DELETE RESTRICT,
    target_id  UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ===========================================================================
-- 6. AI_ENGINE schema tables
-- ===========================================================================

-- t_candidate_scores ---------------------------------------------------------
CREATE TABLE ai_engine.t_candidate_scores (
    candidate_id UUID PRIMARY KEY REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    match_score  INT,
    analyzed_at  TIMESTAMP
);

-- t_score_factors ------------------------------------------------------------
CREATE TABLE ai_engine.t_score_factors (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id UUID NOT NULL REFERENCES ai_engine.t_candidate_scores (candidate_id) ON DELETE CASCADE,
    type         factor_type NOT NULL,
    description  VARCHAR,
    impact       INT
);

-- t_resume_embeddings --------------------------------------------------------
CREATE TABLE ai_engine.t_resume_embeddings (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id      UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    chunk_text   TEXT,
    embedding    vector(1536)
);

-- t_processing_tasks ---------------------------------------------------------
CREATE TABLE ai_engine.t_processing_tasks (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id      VARCHAR UNIQUE,
    entity_id        UUID,
    status           task_status NOT NULL DEFAULT 'pending',
    progress_percent INT DEFAULT 0,
    error_message    TEXT,
    updated_at       TIMESTAMP
);

-- t_hiring_forecasts ---------------------------------------------------------
CREATE TABLE ai_engine.t_hiring_forecasts (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id          UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    forecast_month   DATE,
    predicted_volume INT,
    generated_at     TIMESTAMP
);

-- t_communications -----------------------------------------------------------
CREATE TABLE ai_engine.t_communications (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id        UUID NOT NULL REFERENCES hiring.t_candidates (id) ON DELETE CASCADE,
    generated_by_user_id UUID NOT NULL REFERENCES auth.t_users (id) ON DELETE CASCADE,
    type                email_type NOT NULL,
    content             TEXT,
    sent_at             TIMESTAMP
);

-- t_chat_sessions ------------------------------------------------------------
CREATE TABLE ai_engine.t_chat_sessions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id             UUID NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.t_users (id) ON DELETE CASCADE,
    type                chat_type NOT NULL,
    target_candidate_id UUID REFERENCES hiring.t_candidates (id) ON DELETE SET NULL,
    title               VARCHAR,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP
);

-- t_chat_messages ------------------------------------------------------------
CREATE TABLE ai_engine.t_chat_messages (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id  UUID NOT NULL REFERENCES ai_engine.t_chat_sessions (id) ON DELETE CASCADE,
    role        message_role NOT NULL,
    content     TEXT,
    tokens_used INT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;
