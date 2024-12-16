create table if not EXISTS user_verification_tokens
(
    email      varchar(100) not null,
    token      varchar(100) not null,
    created_at TIMESTAMPTZ not null default current_timestamp,
    updated_at TIMESTAMPTZ not null default current_timestamp,
    primary key (email),
    foreign key (email) references users (email)
);
