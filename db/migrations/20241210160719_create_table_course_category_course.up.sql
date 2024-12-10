create table course_category_course (
    id SERIAL primary key,
    course_id int not null,
    course_category_id int not null,
    created_at TIMESTAMPTZ not null default current_timestamp,
    updated_at TIMESTAMPTZ not null default current_timestamp,
    foreign key (course_id) references courses(id),
    foreign key (course_category_id) references course_categories(id)
);