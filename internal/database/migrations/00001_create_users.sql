-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    crated_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    online BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS user_inventory (
    id SERIAL,
    user_id UUID PRIMARY KEY REFERENCES users(id),
    item_name TEXT NOT NULL,
    quantity INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS user_resources (
    user_id UUID PRIMARY KEY REFERENCES users(id),

    -- General Info
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    stamina_current INTEGER DEFAULT 100,
    stamina_max INTEGER DEFAULT 100,
    
    -- Currencies
    bits BIGINT DEFAULT 1500,
    yen BIGINT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_resources_user_id ON user_resources(user_id);
CREATE INDEX IF NOT EXISTS idx_user_inventory_user_id ON user_inventory(user_id); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_inventory;
DROP TABLE IF EXISTS user_resources;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_user_resources_user_id;
DROP INDEX IF EXISTS idx_user_inventory_user_id;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
