-- Drop Transactions Table (if needed)
do $$
begin
    if exists (select 1 from information_schema.tables where table_name = 'discounts') then
        drop table discounts;
    end if;
    if exists (select 1 from pg_type where typname = 'discounttype') then
        drop type DiscountType;
    end if;
end $$;

-- -- Drop Enum Type PaymentStatus (if needed)
-- do $$
-- begin
    
-- end $$;
