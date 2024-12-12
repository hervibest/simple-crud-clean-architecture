create table course_sections (
    id SERIAL primary key,
    uuid UUID not null,
    course_id int not null,
    title VARCHAR(500) not null unique,
    sequence int ,
    description TEXT,
    foreign key (course_id) references courses(id),
    created_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

-- Add hash index on uuid column
create index idx_course_sections_uuid on course_sections using hash (uuid);