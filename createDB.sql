-- This is SQLite script

-- Stores only the useful details about the users
CREATE TABLE users (
    id INTEGER PRIMARY KEY CHECK(id BETWEEN 0 AND 16777215),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    daily_limit INTEGER NOT NULL DEFAULT 100,
    sent_today INTEGER NOT NULL DEFAULT 0,
    reset_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    emails JSON NOT NULL CHECK(json_valid(emails))
);

-- Keeps track of message/email status
CREATE TABLE email_records (
    message_id INTEGER PRIMARY KEY AUTOINCREMENT CHECK (message_id BETWEEN 0 AND 16777215),
    user_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK(action IN ('IN_QUEUE', 'SUCCESS', 'FAILURE')) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger to update the 'updated_at' column of email_records
CREATE TRIGGER update_email_records_updated_at
    AFTER UPDATE ON email_records
    FOR EACH ROW
BEGIN
    UPDATE email_records
    SET updated_at = CURRENT_TIMESTAMP
    WHERE message_id = OLD.message_id;
END;



---------------Sample data-------------------------------
INSERT INTO users (emails) VALUES
    ('["user1@example.com", "alt1@example.com"]'),
    ('["user2@example.com"]'),
    ('["user3@example.com", "alt3@example.com"]'),
    ('["user4@example.com"]'),
    ('["user5@example.com", "alt5@example.com", "extra5@example.com"]'),
    ('["user6@example.com"]'),
    ('["user7@example.com", "alt7@example.com"]'),
    ('["user8@example.com"]'),
    ('["user9@example.com"]'),
    ('["user10@example.com", "alt10@example.com"]');

INSERT INTO email_records (user_id, action) VALUES
    (1, 'IN_QUEUE'),
    (1, 'SUCCESS'),
    (2, 'FAILURE'),
    (2, 'IN_QUEUE'),
    (3, 'SUCCESS'),
    (3, 'SUCCESS'),
    (4, 'FAILURE'),
    (4, 'IN_QUEUE'),
    (5, 'SUCCESS'),
    (5, 'SUCCESS'),
    (5, 'SUCCESS'),
    (6, 'IN_QUEUE'),
    (6, 'FAILURE'),
    (7, 'SUCCESS'),
    (7, 'IN_QUEUE'),
    (8, 'IN_QUEUE'),
    (9, 'FAILURE'),
    (9, 'SUCCESS'),
    (10, 'SUCCESS'),
    (10, 'IN_QUEUE'),
    (10, 'FAILURE');

UPDATE users
SET sent_today =
    (SELECT COUNT(*) FROM email_records
        WHERE email_records.user_id = users.id AND action = 'SUCCESS');
