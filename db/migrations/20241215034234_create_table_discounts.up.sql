-- Create Enum Type for PaymentStatus
do $$
begin
    if not exists (select 1 from pg_type where typname = 'discounttype') then
        create type DiscountType as enum ('PERCENT', 'REBATE');
    end if;
end $$;

-- Create Transactions Table
create table if not exists discounts (
    id serial primary key,
    uuid UUID not null,
    name varchar(255),
    value float not null,
    is_active bool not null default false,
    type DiscountType not null,
    valid_until TIMESTAMPTZ,
    start_active_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ default current_timestamp not null,
    updated_at TIMESTAMPTZ default current_timestamp not null
);

-- Add hash index on uuid column
create index idx_discounts_uuid on discounts using hash (uuid);