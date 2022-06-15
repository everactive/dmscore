ALTER TABLE device_version
    DROP COLUMN updated_at;

ALTER TABLE device_version
    DROP COLUMN deleted_at;

ALTER TABLE device_version
    DROP COLUMN created_at;