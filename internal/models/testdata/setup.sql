-- Create the snippets table
CREATE TABLE snippets (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
) ENGINE=InnoDB;

CREATE INDEX idx_snippets_created_at ON snippets (created_at);

-- Create the users table
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_uc_email UNIQUE (email)
) ENGINE=InnoDB;

-- Seed with initial data for testing
-- Note: password is 'password' hashed with bcrypt (cost 12)
INSERT INTO users (name, email, hashed_password, created_at) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTj99w8Atv8AUwHeJ5QH.29Pn0zUX1HG0otL6GZuxX1Vn6ObrO72',
    '2026-01-01 10:00:00'
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.',
    '2026-01-01 10:00:00',
    '2027-01-01 10:00:00'
);
