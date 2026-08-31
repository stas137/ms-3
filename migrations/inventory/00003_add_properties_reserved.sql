-- +goose Up
ALTER TABLE parts
    ADD COLUMN properties JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN reserved INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE parts
    DROP COLUMN IF EXISTS reserved,
    DROP COLUMN IF EXISTS properties;