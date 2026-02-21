-- Create the `snippets` table
CREATE TABLE snippets (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
) ENGINE=InnoDB;

-- Add an index on the created_at column
CREATE INDEX idx_snippets_created_at ON snippets(created_at);
