-- Drop index idx_vouchers_uuid if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM pg_indexes 
        WHERE indexname = 'idx_vouchers_uuid'
    ) THEN
        DROP INDEX idx_vouchers_uuid;
    END IF;
END $$;

-- Drop table vouchers if it exists
DROP TABLE IF EXISTS vouchers;

-- Drop type VoucherType if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM pg_type 
        WHERE typname = 'vouchertype'
    ) THEN
        DROP TYPE VoucherType;
    END IF;
END $$;
