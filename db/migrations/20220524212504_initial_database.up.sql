-- 20220524212504_initial_database.up.sql


CREATE TABLE public.openidnonce (
                                    id integer NOT NULL,
                                    nonce character varying(255) NOT NULL,
                                    endpoint character varying(255) NOT NULL,
                                    "timestamp" integer NOT NULL
);

CREATE SEQUENCE public.openidnonce_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



CREATE TABLE public.organization (
                                     id integer NOT NULL,
                                     code character varying(200) NOT NULL,
                                     name character varying(200) NOT NULL,
                                     active boolean DEFAULT true
);



CREATE SEQUENCE public.organization_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.organization_id_seq OWNED BY public.organization.id;


CREATE TABLE public.organization_user (
                                          id integer NOT NULL,
                                          org_id character varying(200) NOT NULL,
                                          username character varying(200) NOT NULL
);

CREATE SEQUENCE public.organization_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.organization_user_id_seq OWNED BY public.organization_user.id;


CREATE TABLE public.settings (
                                 id bigint NOT NULL,
                                 created_at timestamp with time zone,
                                 updated_at timestamp with time zone,
                                 deleted_at timestamp with time zone,
                                 key text,
                                 value text
);


CREATE SEQUENCE public.settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.settings_id_seq OWNED BY public.settings.id;

CREATE TABLE public.userinfo (
                                 id integer NOT NULL,
                                 created timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                 modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                 username character varying(200) NOT NULL,
                                 name character varying(200) NOT NULL,
                                 email character varying(200) NOT NULL,
                                 user_role integer NOT NULL
);

CREATE SEQUENCE public.userinfo_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.userinfo_id_seq OWNED BY public.userinfo.id;



ALTER TABLE ONLY public.openidnonce ALTER COLUMN id SET DEFAULT nextval('public.openidnonce_id_seq'::regclass);


ALTER TABLE ONLY public.organization ALTER COLUMN id SET DEFAULT nextval('public.organization_id_seq'::regclass);


ALTER TABLE ONLY public.organization_user ALTER COLUMN id SET DEFAULT nextval('public.organization_user_id_seq'::regclass);


ALTER TABLE ONLY public.settings ALTER COLUMN id SET DEFAULT nextval('public.settings_id_seq'::regclass);


ALTER TABLE ONLY public.userinfo ALTER COLUMN id SET DEFAULT nextval('public.userinfo_id_seq'::regclass);


ALTER TABLE ONLY public.openidnonce
    ADD CONSTRAINT openidnonce_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_code_key UNIQUE (code);


ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.organization_user
    ADD CONSTRAINT organization_user_org_id_username_key UNIQUE (org_id, username);

ALTER TABLE ONLY public.organization_user
    ADD CONSTRAINT organization_user_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.userinfo
    ADD CONSTRAINT userinfo_pkey PRIMARY KEY (id);


ALTER TABLE ONLY public.userinfo
    ADD CONSTRAINT userinfo_username_key UNIQUE (username);


CREATE INDEX idx_settings_deleted_at ON public.settings USING btree (deleted_at);
