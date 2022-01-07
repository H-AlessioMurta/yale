\connect books
ALTER DATABASE "books" SET TIMEZONE TO 'Europe/Rome';
SET TIMEZONE TO 'Europe/Rome';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "books"
(
    id uuid DEFAULT uuid_generate_v4 (),
    title character varying,
    authors character varying
    
) TABLESPACE pg_default;



ALTER TABLE "books"
    OWNER to postgres;