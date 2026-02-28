-- =============================================================================
-- Migration: 000003_invite_tables (UP)
-- Description: Introduce a dedicated invites table (separating pending invites
--              from active users), plus many-to-many job-access tables for both
--              active users and pending invites (RBAC).
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1. auth.t_invites — pending invitations (not yet accepted)
--    An invite ceases to exist once it is accepted or explicitly cancelled.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS auth.t_invites (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id    UUID        NOT NULL REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    email      VARCHAR     NOT NULL,
    role       user_role   NOT NULL,
    token      VARCHAR     NOT NULL,
    expires_at TIMESTAMP   NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_invites_token UNIQUE (token)
);

-- Fast lookup by email (e.g. "does a pending invite already exist for this email?")
CREATE INDEX IF NOT EXISTS idx_invites_email ON auth.t_invites (email);

-- ---------------------------------------------------------------------------
-- 2. hiring.t_job_access — RBAC: active user ↔ job (many-to-many)
--    Populated when an invite is accepted (transferred from t_invite_job_access).
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS hiring.t_job_access (
    user_id UUID NOT NULL REFERENCES auth.t_users (id)  ON DELETE CASCADE,
    job_id  UUID NOT NULL REFERENCES hiring.t_jobs  (id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, job_id)
);

-- ---------------------------------------------------------------------------
-- 3. auth.t_invite_job_access — RBAC: pending invite ↔ job (many-to-many)
--    Deleted (via CASCADE) once the parent invite row is deleted.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS auth.t_invite_job_access (
    invite_id UUID NOT NULL REFERENCES auth.t_invites  (id) ON DELETE CASCADE,
    job_id    UUID NOT NULL REFERENCES hiring.t_jobs   (id) ON DELETE CASCADE,

    PRIMARY KEY (invite_id, job_id)
);

COMMIT;
