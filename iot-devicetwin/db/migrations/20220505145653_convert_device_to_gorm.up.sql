-- 20220505145653_convert_device_to_gorm.up.sql

ALTER TABLE device
    ADD created_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device
    ADD deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device
    ADD updated_at TIMESTAMP WITH TIME ZONE;

UPDATE device SET updated_at = current_timestamp;

UPDATE device SET created_at = created;

ALTER TABLE device
    DROP COLUMN created;