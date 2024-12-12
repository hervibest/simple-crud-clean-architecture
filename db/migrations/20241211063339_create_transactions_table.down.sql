-- Drop Transactions Table (if needed)
do $$
begin
    if exists (select 1 from information_schema.tables where table_name = 'transactions') then
        drop table transactions;
    end if;
end $$;

-- Drop Enum Type PaymentStatus (if needed)
do $$
begin
    if exists (select 1 from pg_type where typname = 'paymentstatus') then
        drop type PaymentStatus;
    end if;
end $$;
