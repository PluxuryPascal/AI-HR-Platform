BEGIN;

ALTER TABLE auth.t_invites DROP COLUMN status;

DROP TYPE invite_status;

COMMIT;
