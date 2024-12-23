create table if not EXISTS careers (
    id SERIAL primary key,
    uuid UUID not null,
    title VARCHAR(255) not null unique,
    slug VARCHAR(255) not null unique,
    description TEXT,
    is_active boolean not null default true,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_careers_uuid on careers using hash (uuid);