ALTER TABLE action
    ADD created_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE action
    ADD deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE action
    ADD updated_at TIMESTAMP WITH TIME ZONE;

UPDATE action SET updated_at = modified;
UPDATE action SET created_at = created;

ALTER TABLE action
    DROP COLUMN created;

ALTER TABLE action
    DROP COLUMN modified;