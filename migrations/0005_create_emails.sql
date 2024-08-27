CREATE TABLE IF NOT EXISTS emails
(
  address         TEXT NOT NULL,
  verified        INTEGER NOT NULL DEFAULT 0,
  uid             INTEGER NOT NULL,
  id              INTEGER PRIMARY KEY,
  created_at      INTEGER DEFAULT(unixepoch('now')),
  updated_at      INTEGER DEFAULT(unixepoch('now')),
  FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX emails_uid ON emails(uid);
CREATE UNIQUE INDEX emails_uid_address ON emails(uid, address);

CREATE TRIGGER emails_updated_at AFTER UPDATE ON emails
  BEGIN
      UPDATE emails
      SET updated_at = unixepoch('now')
      WHERE id = old.id;
  END;
