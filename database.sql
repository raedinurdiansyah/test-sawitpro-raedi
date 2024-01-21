/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE FUNCTION set_last_modified_at() RETURNS trigger AS $$
BEGIN
    NEW.last_modified_at := NOW()::timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE users (
  "id" serial PRIMARY KEY,
  "guid" UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
  "full_name" VARCHAR (60) NOT NULL,
  "phone_number" VARCHAR (50) NOT NULL,
  "password" VARCHAR (255) NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  "last_modified_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE
);

-- ALTER TABLE users
-- ADD CONSTRAINT users_unique_phone_number_password_key
-- UNIQUE (phone_number, password);

CREATE UNIQUE INDEX users_unique_phone_number_password_key
ON users (phone_number, password)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX users_unique_phone_number_key ON users (phone_number) WHERE deleted_at IS NULL;


CREATE TRIGGER set_last_modified_at BEFORE
UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE set_last_modified_at();