create table if not EXISTS certificate_categories (
    id SERIAL primary key,
    uuid UUID not null,
    name VARCHAR(255) not null unique,
    slug VARCHAR(255) not null unique,
    description TEXT,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_certificate_categories_uuid on certificate_categories using hash (uuid);