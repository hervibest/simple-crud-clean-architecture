ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS voucher_id INT;

ALTER TABLE transactions 
ADD CONSTRAINT fk_voucher
FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;
