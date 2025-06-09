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

-- TODO: Remove dummy data
INSERT INTO mintages_s (
	country,
	type,
	year,
	mintmark,
	[€0,01],
	[€0,02],
	[€0,05],
	[€0,10],
	[€0,20],
	[€0,50],
	[€1,00],
	[€2,00],
	reference
) VALUES
	("ad", 0, 2014, NULL, 60000, 60000, 860000, 860000, 860000, 340000, 511843, 360000, NULL),
	("ad", 0, 2015, NULL, 0, 0, 0, 0, 0, 0, 0, 1072400, NULL),
	("ad", 0, 2016, NULL, 0, 0, 0, 0, 0, 0, 2339200, 0, NULL),
	("ad", 0, 2017, NULL, 2582395, 1515000, 2191421, 1103000, 1213000, 968800, 17000, 794588, NULL),
	("ad", 0, 2018, NULL, 2430000, 2550000, 1800000, 980000, 1014000, 890000, 0, 868000, NULL),
	("ad", 0, 2019, NULL, 2447000, 1727000, 2100000, 1610000, 1570000, 930000, 0, 1058310, NULL),
	("ad", 0, 2020, NULL, 0, 0, 0, 860000, 175000, 740000, 0, 1500000, NULL),
	("ad", 0, 2021, NULL, 200000, 700000, 0, 1400000, 1420000, 600000, 50000, 1474500, NULL),
	("ad", 0, 2022, NULL, 700000, 450000, 400000, 700000, 700000, 380000, 0, 1708000, NULL),
	("ad", 0, 2023, NULL, 0, 0, 0, 0, 0, 0, 0, 2075250, NULL),
	("ad", 0, 2024, NULL, 0, 900300, 1950000, 1000000, 700000, 500000, 1050000, 1601200, NULL),
	("ad", 0, 2025, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);

CREATE TABLE users (
	email    TEXT COLLATE BINARY,
	username TEXT COLLATE BINARY,
	password TEXT COLLATE BINARY,
	adminp   INTEGER
);