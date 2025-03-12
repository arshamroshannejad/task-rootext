CREATE TABLE users (
    id SERIAL PRIMARY KEY,              -- Unique identifier for each user (auto-incrementing)
    email VARCHAR(255) NOT NULL UNIQUE, -- Email address, must be unique
    password VARCHAR(255) NOT NULL,    -- Password (store hashed passwords only)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of user creation
);

-- Create an index on the email column for faster lookups
CREATE INDEX idx_email ON users (email);
