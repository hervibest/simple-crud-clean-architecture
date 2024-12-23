create table if not EXISTS certificates (
    id SERIAL primary key,
    uuid UUID not null,
    categories_id int not null,
    name VARCHAR(255) not null unique,
    slug VARCHAR(255) not null unique,
    description TEXT,
    price double precision,
    is_active boolean not null default true,
    foreign key (categories_id) references certificate_categories(id),
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_certificates_uuid on certificates using hash (uuid);