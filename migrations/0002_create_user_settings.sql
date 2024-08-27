CREATE TABLE IF NOT EXISTS user_settings
(
  theme       TEXT NOT NULL,
  uid         INTEGER NOT NULL,
  id          INTEGER PRIMARY KEY,
  created_at  INTEGER DEFAULT(unixepoch('now')),
  updated_at  INTEGER DEFAULT(unixepoch('now')),
  FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX user_settings_uid ON user_settings(uid);

CREATE TRIGGER user_settings_updated_at AFTER UPDATE ON user_settings
  BEGIN
      UPDATE user_settings
      SET updated_at = unixepoch('now')
      WHERE id = old.id;
  END;
