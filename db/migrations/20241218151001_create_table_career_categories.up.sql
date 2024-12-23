create table if not EXISTS career_category_career (
    id SERIAL primary key,
    career_id int not null,
    career_category_id int not null,
    created_at TIMESTAMPTZ not null default current_timestamp,
    updated_at TIMESTAMPTZ not null default current_timestamp,
    foreign key (career_id) references careers(id),
    foreign key (career_category_id) references career_categories(id)
);