create table employees (
    id SERIAL primary key,
    uuid UUID not null,
    email VARCHAR(100) not null unique,
    password VARCHAR(100) not null,
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_employees_uuid on employees using hash (uuid);