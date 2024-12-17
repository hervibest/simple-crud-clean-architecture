-- Hapus foreign key constraint jika ada
ALTER TABLE transactions 
DROP CONSTRAINT IF EXISTS fk_voucher;

-- Hapus kolom voucher_id jika ada
ALTER TABLE transactions 
DROP COLUMN IF EXISTS voucher_id;
