CREATE TABLE health_hashes (
                        id int generated always as identity,
                        created_at timestamptz,
                        deleted_at timestamptz,
                        updated_at timestamptz,
                        last_refresh timestamptz,
                        org_id character varying(200) NOT NULL,
                        device_id character varying(200) NOT NULL,
                        snap_list_hash character varying(256) NOT NULL,
                        installed_snaps_hash character varying(256) NOT NULL
);