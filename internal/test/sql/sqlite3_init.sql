CREATE TABLE people (
  id integer PRIMARY KEY AUTOINCREMENT,
  name varchar NOT NULL,
  email varchar,
  created_at datetime NOT NULL,
  updated_at datetime
);

CREATE TABLE projects (
  name varchar NOT NULL,
  id varchar PRIMARY KEY,
  start date NOT NULL,
  end date
);

CREATE TABLE person_project (
  person_id integer NOT NULL REFERENCES people ON DELETE CASCADE,
  project_id varchar NOT NULL REFERENCES projects ON DELETE CASCADE,
  UNIQUE (person_id, project_id)
);
