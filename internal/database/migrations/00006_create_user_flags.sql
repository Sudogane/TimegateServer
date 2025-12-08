-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_flags (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    flag_key TEXT NOT NULL UNIQUE,
    flag_value JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, flag_key)
);

CREATE INDEX IF NOT EXISTS idx_user_flags_user_id ON user_flags(user_id);
CREATE INDEX IF NOT EXISTS idx_user_flags_flag_key ON user_flags(flag_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_flags;
DROP INDEX IF EXISTS idx_user_flags_user_id;
DROP INDEX IF EXISTS idx_user_flags_flag_key;
-- +goose StatementEnd
