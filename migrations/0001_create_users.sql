CREATE TABLE IF NOT EXISTS users
(
  pw_hash     TEXT NOT NULL,
  username    TEXT NOT NULL,
  id          INTEGER PRIMARY KEY,
  created_at  INTEGER DEFAULT(unixepoch('now')),
  updated_at  INTEGER DEFAULT(unixepoch('now'))
);

CREATE UNIQUE INDEX users_username ON users(username);

CREATE TRIGGER users_updated_at AFTER UPDATE ON users
  BEGIN
      UPDATE users
      SET updated_at = unixepoch('now')
      WHERE id = old.id;
  END;
