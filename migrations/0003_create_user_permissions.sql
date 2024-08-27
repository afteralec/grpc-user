CREATE TABLE IF NOT EXISTS user_permissions 
(
  name        TEXT NOT NULL,
  iuid        INTEGER NOT NULL,
  uid         INTEGER NOT NULL,
  id          INTEGER PRIMARY KEY,
  created_at  INTEGER DEFAULT(unixepoch('now')),
  FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX user_permissions_name ON user_permissions(uid, name);
CREATE INDEX user_permissions_uid ON user_permissions(uid);
CREATE INDEX user_permissions_iuid ON user_permissions(iuid);
