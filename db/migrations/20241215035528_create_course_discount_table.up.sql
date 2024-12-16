create table if not EXISTS course_discount (
    id SERIAL primary key,
    course_id int not null,
    discount_id int not null,
    created_at TIMESTAMPTZ not null default current_timestamp,
    updated_at TIMESTAMPTZ not null default current_timestamp,
    foreign key (course_id) references courses(id),
    foreign key (discount_id) references discounts(id)
);