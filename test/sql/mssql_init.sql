CREATE TABLE [people] (
  [id] int identity(1, 1) PRIMARY KEY,
  [group_id] int DEFAULT 65534,
  [name] varchar(255) NOT NULL,
  [email] varchar(255),
  [created_at] datetime2 NOT NULL,
  [updated_at] datetime2
);

CREATE TABLE [projects] (
  [name] varchar(255) NOT NULL,
  [id] varchar(255) PRIMARY KEY,
  [start] date NOT NULL,
  [end] date
);

CREATE TABLE [person_project] (
  [person_id] int NOT NULL REFERENCES [people] ON DELETE CASCADE,
  [project_id] varchar(255) NOT NULL REFERENCES [projects] ON DELETE CASCADE,
  UNIQUE ([person_id], [project_id])
);

CREATE TABLE id_only (
  [id] int identity(1, 1) PRIMARY KEY
);

-- to allow insert test data with IDs
SET IDENTITY_INSERT people ON;
