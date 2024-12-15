create table files (
    id SERIAL primary key,
    uuid UUID not null,
    filename VARCHAR(500) not null,
    mimetype VARCHAR(500) not null,
    path VARCHAR(500) not null,
    size int not null,
    fileable_id int not null,
    fileable_type VARCHAR(100) not null,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_files_uuid on files using hash (uuid);