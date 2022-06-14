CREATE TABLE action (
                               id integer NOT NULL,
                               created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               org_id character varying(200) NOT NULL,
                               device_id character varying(200) NOT NULL,
                               action_id character varying(200) NOT NULL,
                               action character varying(200) NOT NULL,
                               status character varying(200) DEFAULT ''::character varying,
                               message text DEFAULT ''::text
);



CREATE SEQUENCE action_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



CREATE TABLE device (
                               id integer NOT NULL,
                               created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               lastrefresh timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               org_id character varying(200) NOT NULL,
                               device_id character varying(200) NOT NULL,
                               brand character varying(200) NOT NULL,
                               model character varying(200) NOT NULL,
                               serial character varying(200) NOT NULL,
                               store_id character varying(200) NOT NULL,
                               device_key text,
                               active boolean DEFAULT true
);



CREATE SEQUENCE device_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



CREATE TABLE device_snap (
                                    id integer NOT NULL,
                                    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                    device_id bigint NOT NULL,
                                    name character varying(200) NOT NULL,
                                    installed_size bigint DEFAULT 0,
                                    installed_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                    status character varying(200) DEFAULT ''::character varying,
                                    channel character varying(200) DEFAULT ''::character varying,
                                    confinement character varying(200) DEFAULT ''::character varying,
                                    version character varying(200) DEFAULT ''::character varying,
                                    revision bigint DEFAULT 0,
                                    devmode boolean DEFAULT false,
                                    config text DEFAULT ''::text
);


CREATE SEQUENCE device_snap_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


CREATE TABLE device_version (
                                       id integer NOT NULL,
                                       device_id integer NOT NULL,
                                       version character varying(200) NOT NULL,
                                       series character varying(200) DEFAULT ''::character varying,
                                       os_id character varying(200) DEFAULT ''::character varying,
                                       os_version_id character varying(200) DEFAULT ''::character varying,
                                       on_classic boolean DEFAULT false,
                                       kernel_version character varying(200) DEFAULT ''::character varying
);


CREATE SEQUENCE device_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE group_device_link (
                                          id integer NOT NULL,
                                          created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                          org_id character varying(200) NOT NULL,
                                          group_id integer NOT NULL,
                                          device_id integer NOT NULL
);


CREATE SEQUENCE group_device_link_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


CREATE TABLE org_group (
                                  id integer NOT NULL,
                                  created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                  modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                  org_id character varying(200) NOT NULL,
                                  name character varying(200) NOT NULL
);


CREATE SEQUENCE org_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE service_statuses (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    device_snap_id bigint,
    name text,
    daemon text,
    enabled boolean,
    active boolean
);

CREATE SEQUENCE service_statuses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE service_statuses_id_seq OWNED BY service_statuses.id;

ALTER TABLE ONLY action
    ADD CONSTRAINT action_pkey PRIMARY KEY (id);


ALTER TABLE ONLY device
    ADD CONSTRAINT device_device_id_key UNIQUE (device_id);


ALTER TABLE ONLY device
    ADD CONSTRAINT device_pkey PRIMARY KEY (id);


ALTER TABLE ONLY device_snap
    ADD CONSTRAINT device_snap_pkey PRIMARY KEY (id);


ALTER TABLE ONLY device_version
    ADD CONSTRAINT device_version_device_id_key UNIQUE (device_id);


ALTER TABLE ONLY device_version
    ADD CONSTRAINT device_version_pkey PRIMARY KEY (id);


ALTER TABLE ONLY group_device_link
    ADD CONSTRAINT group_device_link_pkey PRIMARY KEY (id);

ALTER TABLE ONLY org_group
    ADD CONSTRAINT org_group_pkey PRIMARY KEY (id);


ALTER TABLE ONLY service_statuses
    ADD CONSTRAINT service_statuses_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX device_snap_idx ON device_snap USING btree (device_id, name);


CREATE INDEX idx_service_statuses_deleted_at ON service_statuses USING btree (deleted_at);


CREATE INDEX org_group_idx ON org_group USING btree (org_id, name);

ALTER TABLE ONLY device_snap
    ADD CONSTRAINT device_snap_device_id_fkey FOREIGN KEY (device_id) REFERENCES device(id) ON DELETE CASCADE;

ALTER TABLE ONLY device_version
    ADD CONSTRAINT device_version_device_id_fkey FOREIGN KEY (device_id) REFERENCES device(id);


ALTER TABLE ONLY group_device_link
    ADD CONSTRAINT group_device_link_device_id_fkey FOREIGN KEY (device_id) REFERENCES device(id);

ALTER TABLE ONLY group_device_link
    ADD CONSTRAINT group_device_link_group_id_fkey FOREIGN KEY (group_id) REFERENCES org_group(id);

ALTER TABLE ONLY action ALTER COLUMN id SET DEFAULT nextval('action_id_seq'::regclass);


ALTER TABLE ONLY device ALTER COLUMN id SET DEFAULT nextval('device_id_seq'::regclass);


ALTER TABLE ONLY device_snap ALTER COLUMN id SET DEFAULT nextval('device_snap_id_seq'::regclass);


ALTER TABLE ONLY device_version ALTER COLUMN id SET DEFAULT nextval('device_version_id_seq'::regclass);


ALTER TABLE ONLY group_device_link ALTER COLUMN id SET DEFAULT nextval('group_device_link_id_seq'::regclass);


ALTER TABLE ONLY org_group ALTER COLUMN id SET DEFAULT nextval('org_group_id_seq'::regclass);

ALTER TABLE ONLY service_statuses ALTER COLUMN id SET DEFAULT nextval('public.service_statuses_id_seq'::regclass);

