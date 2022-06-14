-- 20220519191940_convert_organization_to_gorm.up.sql

ALTER TABLE organization
    ADD created_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE organization
    ADD deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE organization
    ADD updated_at TIMESTAMP WITH TIME ZONE;

UPDATE organization SET updated_at = current_timestamp;
UPDATE organization SET created_at = current_timestamp;