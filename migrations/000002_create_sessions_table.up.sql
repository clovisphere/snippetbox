-- Create the sessions table
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
) ENGINE=InnoDB;

-- Add an index on the expiry column using your naming convention
CREATE INDEX idx_sessions_expiry ON sessions (expiry);
