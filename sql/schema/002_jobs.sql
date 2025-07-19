-- +goose Up
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    type TEXT UNIQUE NOT NULL,
    salary DECIMAL(10,2) NOT NULL,
    employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE jobs;

