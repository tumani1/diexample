-- +migrate Up
CREATE extension pg_trgm;

CREATE TABLE IF NOT EXISTS fleets
(
    id   BIGSERIAL    NOT NULL,
    name VARCHAR(256) NOT NULL,

    PRIMARY KEY (id)
--     CONSTRAINT uk_fleets_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS calendar
(
    boat_id   BIGINT  NOT NULL,
--         CONSTRAINT calendar_boat_id_fkey REFERENCES "boats" on update cascade on delete restrict,
    date_from date    NOT NULL,
    date_to   date    NOT NULL,
    available BOOLEAN NOT NULL

--     CONSTRAINT uk_calendar_boat_id_date_from UNIQUE (boat_id, date_from, available)
);

CREATE INDEX ix_calendar_boat_id_date_from_date_to ON calendar USING btree ("boat_id", "date_from", "date_to");

CREATE TABLE IF NOT EXISTS yacht_builders
(
    id            BIGSERIAL    NOT NULL,
    name          VARCHAR(256) NOT NULL,

    PRIMARY KEY (id)
--     CONSTRAINT uk_yacht_builders_name UNIQUE (name)
);

CREATE INDEX ix_gin_yacht_builders_name ON yacht_builders USING gin ("name" gin_trgm_ops);
CREATE INDEX ix_btree_yacht_builders_name ON yacht_builders USING btree ("name");

CREATE TABLE IF NOT EXISTS yacht_models
(
    id            BIGSERIAL    NOT NULL,
    name          VARCHAR(256) NOT NULL,
    builder_id    BIGINT       NOT NULL,

    PRIMARY KEY (id)
--     CONSTRAINT uk_yacht_models_name_builder_id UNIQUE (name, builder_id)
);


CREATE INDEX ix_gin_yacht_models_name ON yacht_models USING gin ("name" gin_trgm_ops);
CREATE INDEX ix_btree_yacht_models_builder_id ON yacht_models USING btree (builder_id);
CREATE INDEX ix_btree_yacht_models_name ON yacht_models USING btree ("name");

CREATE TABLE IF NOT EXISTS boats
(
    id       BIGSERIAL NOT NULL,
    fleet_id BIGINT    NOT NULL,
    model_id BIGINT ,

    PRIMARY KEY (id)
--     CONSTRAINT uk_boats_fleet_id_model_id UNIQUE (fleet_id, model_id)
);

CREATE INDEX ix_boats_fleet_id_model_id ON boats USING btree ("fleet_id", "model_id");

-- +migrate Down
DROP TABLE IF EXISTS fleets;
DROP TABLE IF EXISTS yacht_builders;
DROP TABLE IF EXISTS yacht_models;
DROP TABLE IF EXISTS boats;
DROP TABLE IF EXISTS calendar;
