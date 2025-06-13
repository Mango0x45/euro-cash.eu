PRAGMA encoding = "UTF-8";

CREATE TABLE migration (
	id     INTEGER PRIMARY KEY CHECK (id = 1),
	latest INTEGER
);
INSERT INTO migration (id, latest) VALUES (1, -1);

CREATE TABLE mintages_s (
	country   CHAR(2) NOT NULL COLLATE BINARY
		CHECK(length(country) = 2),
	type      INTEGER NOT NULL  -- Codes correspond to contants in mintages.go
		CHECK(type BETWEEN 0 AND 2),
	year      INTEGER NOT NULL,
	mintmark  TEXT,
	[€0,01]   INTEGER,
	[€0,02]   INTEGER,
	[€0,05]   INTEGER,
	[€0,10]   INTEGER,
	[€0,20]   INTEGER,
	[€0,50]   INTEGER,
	[€1,00]   INTEGER,
	[€2,00]   INTEGER,
	reference TEXT
);

CREATE TABLE mintages_c (
	country   CHAR(2) NOT NULL COLLATE BINARY
		CHECK(length(country) = 2),
	type      INTEGER NOT NULL  -- Codes correspond to contants in mintages.go
		CHECK(type BETWEEN 0 AND 2),
	name      TEXT NOT NULL,
	year      INTEGER NOT NULL,
	number    INTEGER NOT NULL,
	mintmark  TEXT,
	mintage   INTEGER,
	reference TEXT
);

CREATE TABLE users (
	email      TEXT COLLATE BINARY,
	username   TEXT COLLATE BINARY,
	password   TEXT COLLATE BINARY,
	adminp     INTEGER,
	translates TEXT COLLATE BINARY
);