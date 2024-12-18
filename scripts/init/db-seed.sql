-- Create the database if it doesn't exist
CREATE DATABASE db;


-- Connect to the newly created or existing database
\c db;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY (START WITH 1),
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    salt VARCHAR NOT NULL
);

-- Insert the default admin user
INSERT INTO users (id, username, password, salt) VALUES
-- password is set to "admin" by default
(0, 'admin', '$argon2id$v=19$m=65536,t=3,p=4$N9Jb/TWjt61CqDLaQoHXHQ$bpiJEOeuQasr/kPGTjZ8v1+EG+dfGzoA3jhauutR+yk', '37d25bfd35a3b7ad42a832da4281d71d')
ON CONFLICT DO NOTHING;