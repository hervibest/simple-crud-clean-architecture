create table users (
    id SERIAL primary key,
    uuid UUID not null,
    email VARCHAR(100) not null unique,
    password VARCHAR(100) not null,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    verified_at TIMESTAMPTZ
);

-- Add hash index on uuid column
create index idx_users_uuid on users using hash (uuid);