-- +migrate Up
COPY boats (id, fleet_id, model_id)
    FROM '/tmp/boats.csv' DELIMITER ',' CSV HEADER;

COPY fleets (id, name)
    FROM '/tmp/fleets.csv' DELIMITER ',' CSV HEADER;

COPY calendar (boat_id, date_from, date_to, available)
    FROM '/tmp/calendar.csv' DELIMITER ',' CSV HEADER;

COPY yacht_builders (id, name)
    FROM '/tmp/yacht_builders.csv' DELIMITER ',' CSV HEADER;

COPY yacht_models (id, name, builder_id)
    FROM '/tmp/yacht_models.csv' DELIMITER ',' CSV HEADER;

-- +migrate Down
TRUNCATE boats;
TRUNCATE fleets;
TRUNCATE calendar;
TRUNCATE yacht_builders;
TRUNCATE yacht_models;
