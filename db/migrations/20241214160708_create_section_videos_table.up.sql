create table section_videos (
    id SERIAL primary key,
    uuid UUID not null,
    section_id int not null,
    title VARCHAR(500) not null unique,
    sequence int ,
    notes TEXT,
    original_name VARCHAR(500),
    original_size double precision,
    original_mime VARCHAR(500),
    media_id VARCHAR(500),
    bucket VARCHAR(500),
    dir VARCHAR(500),
    foreign key (section_id) references course_sections(id),
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_section_videos_uuid on section_videos using hash (uuid);