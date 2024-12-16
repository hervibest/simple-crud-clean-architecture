create table if not EXISTS course_user (
    id SERIAL primary key,
    course_id int not null,
    user_id int not null,
    created_at TIMESTAMPTZ not null default current_timestamp,
    updated_at TIMESTAMPTZ not null default current_timestamp,
    foreign key (course_id) references courses(id),
    foreign key (user_id) references users(id)
);