-- 20220429232505_add_fk_device_snap_service_statuses.down.sql

ALTER TABLE ONLY service_statuses
    DROP CONSTRAINT IF EXISTS fk_device_snap_service_statuses;