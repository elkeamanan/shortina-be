CREATE TABLE links (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL,
    redirection TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

create trigger links_updated_at
    before update
    on links
    for each row
execute procedure update_modified_column();