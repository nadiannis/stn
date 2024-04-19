CREATE TABLE IF NOT EXISTS links (
	id VARCHAR(255) NOT NULL PRIMARY KEY,
	url VARCHAR(2083) NOT NULL,
	back_half VARCHAR(100) NOT NULL,
	engagements BIGINT NOT NULL DEFAULT 0,
	user_id VARCHAR(255) NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

ALTER TABLE links ADD CONSTRAINT links_uc_back_half UNIQUE (back_half);
ALTER TABLE links ADD CONSTRAINT links_fk_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
