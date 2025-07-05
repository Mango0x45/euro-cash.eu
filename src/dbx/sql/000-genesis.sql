PRAGMA encoding = "UTF-8";

CREATE TABLE migration (
	id     INTEGER PRIMARY KEY CHECK (id = 1),
	latest INTEGER
);
INSERT INTO migration (id, latest) VALUES (1, -1);

CREATE TABLE mintages_s (
	country      CHAR(2) NOT NULL COLLATE BINARY
		CHECK(length(country) = 2),
	-- Codes correspond to contants in mintages.go
	type         INTEGER NOT NULL
		CHECK(type BETWEEN 0 AND 2),
	year         INTEGER NOT NULL,
	denomination REAL NOT NULL,
	mintmark     TEXT,
	mintage      INTEGER,
	reference    TEXT
);

CREATE TABLE mintages_c (
	country   CHAR(2) NOT NULL COLLATE BINARY
		CHECK(length(country) = 2),
	-- Codes correspond to contants in mintages.go
	type      INTEGER NOT NULL
		CHECK(type BETWEEN 0 AND 2),
	year      INTEGER NOT NULL,
	name      TEXT NOT NULL,
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