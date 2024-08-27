CREATE TABLE IF NOT EXISTS user_permission_grants
(
  name        TEXT NOT NULL,
  iuid        INTEGER NOT NULL,
  uid         INTEGER NOT NULL,
  id          INTEGER PRIMARY KEY,
  created_at  INTEGER DEFAULT(unixepoch('now'))
);

CREATE INDEX user_permission_grants_uid ON user_permission_grants(uid);
CREATE INDEX user_permission_grants_iuid ON user_permission_grants(iuid);

CREATE TABLE IF NOT EXISTS user_permission_revocations
(
  name        TEXT NOT NULL,
  iuid        INTEGER NOT NULL,
  uid         INTEGER NOT NULL,
  id          INTEGER PRIMARY KEY,
  created_at  INTEGER DEFAULT(unixepoch('now'))
);

CREATE INDEX user_permission_revocations_uid ON user_permission_revocations(uid);
CREATE INDEX user_permission_revocations_iuid ON user_permission_revocations(iuid);
