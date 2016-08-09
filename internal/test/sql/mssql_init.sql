CREATE TABLE [people] (
  [id] int identity(1, 1) PRIMARY KEY,
  [group_id] int DEFAULT 65534,
  [name] varchar(255),
  [email] varchar(255) NULL,
  [created_at] datetime2,
  [updated_at] datetime2 NULL
);

CREATE TABLE [projects] (
  [id] varchar(255) PRIMARY KEY,
  [name] varchar(255),
  [start] date,
  [end] date NULL
);

CREATE TABLE [person_project] (
  [person_id] int REFERENCES [people] ON DELETE CASCADE,
  [project_id] varchar(255) REFERENCES [projects] ON DELETE CASCADE,
  UNIQUE ([person_id], [project_id])
);

CREATE TABLE id_only (
  [id] int identity(1, 1) PRIMARY KEY
);
