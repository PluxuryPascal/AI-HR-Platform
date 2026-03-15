BEGIN;

DROP INDEX IF EXISTS auth.idx_invites_recovery;

ALTER TABLE auth.t_invites 
DROP COLUMN IF EXISTS updated_at;

COMMIT;
