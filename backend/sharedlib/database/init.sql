CREATE TABLE IF NOT EXISTS festival (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		last_used_at BIGINT NOT NULL,
		created_at BIGINT NOT NULL,
		pin VARCHAR(4) NOT NULL DEFAULT 'NONE',
		code VARCHAR(6) NOT NULL UNIQUE,
		password TEXT NOT NULL DEFAULT ''
		);

CREATE TABLE IF NOT EXISTS event (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		active BOOLEAN NOT NULL,
		created_at BIGINT NOT NULL,
		last_used_at BIGINT NOT NULL,
		festival_id BIGINT NOT NULL,
		total INTEGER,
		FOREIGN KEY (festival_id) REFERENCES festival(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS active (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		value INTEGER NOT NULL,
		time BIGINT NOT NULL,
		event_id BIGINT NOT NULL,
		FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS archive (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		value INTEGER NOT NULL,
		time BIGINT NOT NULL,
		event_id BIGINT NOT NULL,
		FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS gauge_max (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		gauge_max INTEGER NOT NULL,
		time BIGINT NOT NULL,
		event_id BIGINT NOT NULL,
		FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS app_user (
		id BIGSERIAL NOT NULL PRIMARY KEY
		);

CREATE TABLE IF NOT EXISTS refresh_token (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		token TEXT NOT NULL,
		expires_at BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		revoked BOOLEAN NOT NULL,
		FOREIGN KEY (user_id) REFERENCES app_user(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS access_token (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		token TEXT NOT NULL,
		expires_at BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		revoked BOOLEAN NOT NULL,
		FOREIGN KEY (user_id) REFERENCES app_user(id) ON DELETE CASCADE
		);

CREATE TABLE IF NOT EXISTS festival_access (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		last_used_at BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		festival_id BIGINT NOT NULL,
		revoked BOOLEAN NOT NULL,
		FOREIGN KEY (user_id) REFERENCES app_user(id) ON DELETE CASCADE,
		FOREIGN KEY (festival_id) REFERENCES festival(id) ON DELETE CASCADE
		);

CREATE INDEX IF NOT EXISTS idx_festival_code ON festival (code);
CREATE INDEX IF NOT EXISTS idx_festival_last_used_at ON festival (last_used_at);

CREATE INDEX IF NOT EXISTS idx_event_festival_id ON event (festival_id);
CREATE INDEX IF NOT EXISTS idx_event_active ON event (active);
CREATE INDEX IF NOT EXISTS idx_event_last_used_at ON event (last_used_at);

CREATE INDEX IF NOT EXISTS idx_active_event_id ON active (event_id);

CREATE INDEX IF NOT EXISTS idx_archived_event_id ON archive (event_id);

CREATE INDEX IF NOT EXISTS idx_gauge_max_event_id_time ON gauge_max (event_id, time DESC);

CREATE INDEX IF NOT EXISTS idx_refresh_token_user_id ON refresh_token (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON refresh_token (token);
CREATE INDEX IF NOT EXISTS idx_refresh_token_expires_at ON refresh_token (expires_at);

CREATE INDEX IF NOT EXISTS idx_access_token_user_id ON access_token (user_id);
CREATE INDEX IF NOT EXISTS idx_access_token_token ON access_token (token);
CREATE INDEX IF NOT EXISTS idx_access_token_expires_at ON access_token (expires_at);

CREATE INDEX IF NOT EXISTS idx_festival_access_festival_id ON festival_access (festival_id);
CREATE INDEX IF NOT EXISTS idx_festival_access_user_id ON festival_access (user_id);
CREATE INDEX IF NOT EXISTS idx_festival_access_last_used_at ON festival_access (last_used_at);
