-- 20220429232505_add_fk_device_snap_service_statuses.up.sql

ALTER TABLE ONLY service_statuses
    ADD CONSTRAINT fk_device_snap_service_statuses FOREIGN KEY (device_snap_id) REFERENCES device_snap(id) ON DELETE CASCADE;