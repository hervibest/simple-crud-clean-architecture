create table if not EXISTS courses (
    id SERIAL primary key,
    uuid UUID not null,
    name VARCHAR(255) not null unique,
    slug VARCHAR(255) not null unique,
    description TEXT,
    price double precision,
    is_active boolean not null default true,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_courses_uuid on courses using hash (uuid);