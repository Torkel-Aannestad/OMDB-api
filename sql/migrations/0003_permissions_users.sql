-- +goose Up
CREATE TABLE IF NOT EXISTS permissions (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  
    code text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES 
    ('movies:read'),
    ('movies:write'),
    ('people:read'), 
    ('people:write'),
    ('casts:read'),
    ('casts:write'),
    ('jobs:read'),
    ('jobs:write'),
    ('categories:read'),
    ('categories:write'),
    ('category-items:read'),
    ('category-items:write'),
    ('movie-links:read'),
    ('movie-links:write'),
    ('people-links:read'),
    ('people-links:write'),
    ('trailers:read'),
    ('trailers:write')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS users_permissions;
DROP TABLE IF EXISTS permissions;