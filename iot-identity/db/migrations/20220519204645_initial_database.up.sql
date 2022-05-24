CREATE TABLE IF NOT EXISTS organization (
    id               serial primary key not null,
    org_id           varchar(200) not null unique,
    name             varchar(200) not null,
    country_name     varchar(200) default '',
    root_cert         text not null,
    root_key          text not null,
    UNIQUE (org_id)
);

CREATE TABLE IF NOT EXISTS device (
    id                serial primary key not null,
    device_id         varchar(200) not null unique,
    org_id            varchar(200) not null,
    brand             varchar(200) not null,
    model             varchar(200) not null,
    serial_number     varchar(200) not null,
    cred_key          text not null,
    cred_cert         text not null,
    cred_mqtt         varchar(200) not null,
    cred_port         varchar(200) not null,

    store_id          varchar(200) default '',
    device_key        text default '',
    status            int default 1,
    device_data       text default '',

    UNIQUE (device_id),
    UNIQUE (brand, model, serial_number)
);

CREATE INDEX IF NOT EXISTS device_id_idx ON device (device_id);

CREATE INDEX IF NOT EXISTS bms_idx ON device (brand, model, serial_number);
