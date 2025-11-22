CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
    updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Trigger to update the updated_at timestamp on record update TODO
-- CREATE TRIGGER update_user_updated_at
-- BEFORE UPDATE ON users
-- FOR EACH ROW
-- BEGIN
--     UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
-- END;