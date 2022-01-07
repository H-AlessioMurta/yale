\connect customers
ALTER DATABASE "customers" SET TIMEZONE TO 'Europe/Rome';
SET TIMEZONE TO 'Europe/Rome';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "customers"
(
    id uuid DEFAULT uuid_generate_v4 (),
    name character varying,
    surname character varying,
    nin character varying

) TABLESPACE pg_default;



ALTER TABLE "customers"
    OWNER to postgres;

