CREATE TABLE [people] (
  [id] int identity(1, 1) PRIMARY KEY,
  [name] varchar(50),
  [email] varchar(50) NULL,
  [created_at] datetime2,
  [updated_at] datetime2 NULL
);

CREATE TABLE [projects] (
  [id] varchar(50) PRIMARY KEY,
  [name] varchar(50),
  [start] date,
  [end] date NULL
);

CREATE TABLE [person_project] (
  [person_id] int REFERENCES [people] ON DELETE CASCADE,
  [project_id] varchar(50) REFERENCES [projects] ON DELETE CASCADE,
  UNIQUE (person_id, project_id)
);
