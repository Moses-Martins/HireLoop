-- +goose Up
CREATE TABLE applications (
    id UUID PRIMARY KEY,
    applicant_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    resume_url TEXT NOT NULL,
    status TEXT NOT NULL
);

-- +goose Down
DROP TABLE applications;
