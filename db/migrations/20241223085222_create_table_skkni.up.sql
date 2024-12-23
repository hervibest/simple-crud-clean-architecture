create table if not EXISTS skkni (
    id SERIAL primary key,
    uuid UUID not null,
    certificate_id int not null,
    name VARCHAR(500) not null unique,
    description TEXT,
    foreign key (certificate_id) references certificates(id),
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_skkni_uuid on skkni using hash (uuid);