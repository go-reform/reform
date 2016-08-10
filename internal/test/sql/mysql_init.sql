SET GLOBAL sql_mode = 'ANSI';

CREATE TABLE people (
  id int NOT NULL AUTO_INCREMENT,
  group_id int DEFAULT 65534,
  name varchar(255) NOT NULL,
  email varchar(255),
  created_at datetime NOT NULL,
  updated_at datetime,
  -- TODO doesn't work, file an issue for github.com/go-sql-driver/mysql
  -- created_at timestamp NOT NULL,
  -- updated_at timestamp,
  PRIMARY KEY (id)
);

CREATE TABLE projects (
  name varchar(255) NOT NULL,
  id varchar(255) NOT NULL,
  start date NOT NULL,
  end date,
  PRIMARY KEY (id)
);

-- https://dev.mysql.com/doc/refman/5.7/en/create-table.html
-- MySQL parses but ignores “inline REFERENCES specifications” (as defined in the SQL standard)
-- where the references are defined as part of the column specification. MySQL accepts REFERENCES
-- clauses only when specified as part of a separate FOREIGN KEY specification.

CREATE TABLE person_project (
  person_id int NOT NULL,
  project_id varchar(255) NOT NULL,
  UNIQUE (person_id, project_id),
  FOREIGN KEY (person_id) REFERENCES people (id) ON DELETE CASCADE,
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TABLE id_only (
  id int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (id)
);
