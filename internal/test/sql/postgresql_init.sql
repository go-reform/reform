CREATE TABLE people (
  id serial PRIMARY KEY,
  group_id integer DEFAULT 65534,
  name varchar NOT NULL,
  email varchar,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone
  -- created_at timestamp without time zone NOT NULL,
  -- updated_at timestamp without time zone
);

CREATE TABLE projects (
  name varchar NOT NULL,
  id varchar PRIMARY KEY,
  start date NOT NULL,
  "end" date
);

CREATE TABLE person_project (
  person_id integer NOT NULL REFERENCES people ON DELETE CASCADE,
  project_id varchar NOT NULL REFERENCES projects ON DELETE CASCADE,
  UNIQUE (person_id, project_id)
);

CREATE TABLE id_only (
  id serial PRIMARY KEY
);

CREATE SCHEMA legacy;

CREATE TABLE legacy.people (
  id serial PRIMARY KEY,
  name varchar
);
