CREATE TABLE gauges (
    name varchar(128) PRIMARY KEY,
    value double precision
);

CREATE TABLE counters (
    name varchar(128) PRIMARY KEY,
    value bigint
);