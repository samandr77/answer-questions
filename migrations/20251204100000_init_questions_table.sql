-- +goose Up
-- Create questions table
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    text VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
-- Drop questions table
DROP TABLE questions;
