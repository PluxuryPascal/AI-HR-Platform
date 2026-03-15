-- =============================================================================
-- Migration: 000010_candidate_parsing_status (UP)
-- Description: Add parsing_status ENUM and column to t_candidates
-- =============================================================================

BEGIN;

CREATE TYPE candidate_parsing_status AS ENUM (
    'pending',       -- кандидат создан, резюме ещё не загружено или в очереди
    'processing',    -- файл загружен, Temporal workflow запущен
    'needs_review',  -- AI распарсил частично, рекрутер должен дополнить данные вручную
    'completed',     -- AI успешно распарсил, данные заполнены
    'failed'         -- OCR провалился / TTL истёк / DLQ — резюме не обработано
);

ALTER TABLE hiring.t_candidates
    ADD COLUMN parsing_status candidate_parsing_status NOT NULL DEFAULT 'pending';

COMMIT;
