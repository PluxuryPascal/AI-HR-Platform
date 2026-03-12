BEGIN;

CREATE TYPE invite_status AS ENUM ('pending', 'accepted', 'processing', 'completed', 'failed');

ALTER TABLE auth.t_invites 
ADD COLUMN status invite_status NOT NULL DEFAULT 'pending';

COMMIT;
