BEGIN;

ALTER TABLE auth.t_invites 
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- Индекс для быстрого поиска зависших инвайтов в статусе processing
CREATE INDEX idx_invites_recovery ON auth.t_invites (status, updated_at) 
WHERE status = 'processing';

COMMIT;
