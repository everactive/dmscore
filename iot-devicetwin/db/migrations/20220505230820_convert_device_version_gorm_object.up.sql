ALTER TABLE device_version
    ADD created_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device_version
    ADD deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device_version
    ADD updated_at TIMESTAMP WITH TIME ZONE;

UPDATE device_version SET updated_at = current_timestamp;
UPDATE device_version SET created_at = current_timestamp;
