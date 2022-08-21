-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklist (
    id SERIAL PRIMARY KEY,
    subnet INET NOT NULL UNIQUE
);

CREATE TABLE whitelist (
    id SERIAL PRIMARY KEY,
    subnet INET NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blacklist;
DROP TABLE IF EXISTS whitelist;
-- +goose StatementEnd
