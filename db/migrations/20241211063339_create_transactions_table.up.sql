-- Create Enum Type for PaymentStatus
do $$
begin
    if not exists (select 1 from pg_type where typname = 'paymentstatus') then
        create type PaymentStatus as enum ('SUCCESS', 'FAILED', 'PENDING', 'EXPIRED', 'CANCELLED', 'REFUND');
    end if;
end $$;

-- Create Transactions Table
create table if not exists transactions (
    id serial primary key,
    trx_id UUID not null unique,
    user_id int not null references users(id) on delete cascade,
    course_id int not null references courses(id) on delete cascade,
    amount double precision not null,
    status PaymentStatus not null,
    snap_token varchar(500),
    external_status varchar(255),
    external_callback_response json,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ default current_timestamp not null,
    updated_at TIMESTAMPTZ default current_timestamp not null
);