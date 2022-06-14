-- 20220505145653_convert_device_to_gorm.down.sql

ALTER TABLE device
    DROP COLUMN updated_at;

ALTER TABLE device
    DROP COLUMN deleted_at;

ALTER TABLE device
    ADD created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP;

UPDATE device set created = created_at;

ALTER TABLE device
    DROP COLUMN created_at;
