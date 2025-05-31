ALTER TABLE links
    DROP CONSTRAINT unique_key,
    DROP column status,
    DROP column created_by;

DROP TYPE IF EXISTS link_status;