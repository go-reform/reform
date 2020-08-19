INSERT INTO people (id, name, email, created_at) VALUES (1, 'Denis Mills', NULL, '2009-11-10 23:00:00');
INSERT INTO people (id, name, email, created_at) VALUES (2, 'Garrick Muller', 'muller_garrick@example.com', '2009-12-12 12:34:56');

INSERT INTO people (id, name, email, created_at) VALUES (101, 'Noble Schumm', NULL, '2013-01-01 00:00:00');
INSERT INTO people (id, name, email, created_at) VALUES (102, 'Elfrieda Abbott', 'elfrieda_abbott@example.org', '2014-01-01 00:00:00');
INSERT INTO people (id, name, email, created_at) VALUES (103, 'Elfrieda Abbott', NULL, '2014-01-01 00:00:00');

-- ANSI quotes for keyword "end"
INSERT INTO projects (i, id, name, start, "end") VALUES (1, 'baron', 'Vicious Baron', '2014-06-01', '2016-02-21');
INSERT INTO projects (i, id, name, start, "end") VALUES (2, 'queen', 'Thirsty Queen', '2016-01-15', NULL);
INSERT INTO projects (i, id, name, start, "end") VALUES (3, 'traveler', 'Kosher Traveler', '2016-02-01', NULL);
INSERT INTO projects (i, id, name, start, "end") VALUES (4, 'lightfoot', 'Sweet Lightfoot', '2016-01-01', NULL);
INSERT INTO projects (i, id, name, start, "end") VALUES (5, 'walker', 'Eager Walker', '2015-01-01', NULL);

INSERT INTO person_project (project_id, person_id) VALUES ('baron', 101);
INSERT INTO person_project (project_id, person_id) VALUES ('baron', 102);
INSERT INTO person_project (project_id, person_id) VALUES ('baron', 103);

INSERT INTO person_project (project_id, person_id) VALUES ('queen', 102);
INSERT INTO person_project (project_id, person_id) VALUES ('queen', 103);

INSERT INTO person_project (project_id, person_id) VALUES ('traveler', 103);
