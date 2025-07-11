GRANT ALL PRIVILEGES ON DATABASE auth_testcase TO testcase_api;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c auth_testcase

CREATE TABLE IF NOT EXISTS users (
    guid uuid NOT NULL DEFAULT uuid_generate_v4(),
    first_name character varying(25) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(25) COLLATE pg_catalog."default" NOT NULL,
    email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (guid),
    CONSTRAINT uni_users_email UNIQUE (email)
);


ALTER TABLE IF EXISTS users
    OWNER to testcase_api;

CREATE TABLE IF NOT EXISTS tokens (
    id bigint PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_guid uuid NOT NULL,
    refresh_token character varying(255) COLLATE pg_catalog."default",
    last_ip character varying(255) COLLATE pg_catalog."default",
    last_fingerprint character varying(255) COLLATE pg_catalog."default",
    pair_id character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT uni_tokens_pair_id UNIQUE (pair_id),
    CONSTRAINT uni_tokens_refresh_token UNIQUE (refresh_token),
    CONSTRAINT uni_tokens_user_guid UNIQUE (user_guid),
    CONSTRAINT fk_tokens_user FOREIGN KEY (user_guid)
        REFERENCES public.users (guid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

ALTER TABLE IF EXISTS tokens
    OWNER to testcase_api;

COPY users(guid,first_name,last_name,email) FROM '/docker-entrypoint-initdb.d/data.csv' DELIMITER ',' CSV HEADER;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testcase_api;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testcase_api;