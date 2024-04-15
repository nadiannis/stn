CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(255) NOT NULL PRIMARY KEY,
	email VARCHAR(254) NOT NULL,
	password CHAR(60) NOT NULL,
	date_joined DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);