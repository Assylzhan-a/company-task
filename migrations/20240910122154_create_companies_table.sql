-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS companies (
                           id UUID PRIMARY KEY,
                           name VARCHAR(15) NOT NULL UNIQUE,
                           description TEXT,
                           amount_of_employees INTEGER NOT NULL,
                           registered BOOLEAN NOT NULL,
                           type VARCHAR(20) NOT NULL,
                           created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                           updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS companies;
