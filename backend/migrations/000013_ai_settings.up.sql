CREATE TABLE ai_engine.t_team_settings (
    team_id       UUID PRIMARY KEY REFERENCES auth.t_teams (id) ON DELETE CASCADE,
    api_key       VARCHAR,
    parse_model   VARCHAR,
    score_model   VARCHAR,
    embed_model   VARCHAR,
    chat_model    VARCHAR,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
