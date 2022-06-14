-- 20220505230850_convert_device_snap_gorm_object

ALTER TABLE device_snap
    ADD created_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device_snap
    ADD deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE device_snap
    ADD updated_at TIMESTAMP WITH TIME ZONE;

UPDATE device_snap SET updated_at = modified;
UPDATE device_snap SET created_at = created;

ALTER TABLE device_snap
    DROP COLUMN created;

ALTER TABLE device_snap
    DROP COLUMN modified;