CREATE TYPE link_status AS ENUM ('unspecified','active', 'inactive');

ALTER TABLE links
    ADD CONSTRAINT unique_key UNIQUE (key),
    ADD COLUMN status link_status NOT NULL DEFAULT 'unspecified',
    ADD COLUMN created_by UUID REFERENCES users(id);