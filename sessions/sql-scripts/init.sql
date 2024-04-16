DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id                  serial               PRIMARY KEY,
    first_name          varchar(20)          NOT NULL,
    last_name           varchar(20)          NOT NULL,
    email               varchar(60)          UNIQUE NOT NULL,
    password            varchar(72)          NOT NULL,
    jwts                varchar(72)[]        NOT NULL
);