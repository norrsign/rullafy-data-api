CREATE TABLE users (
  id       TEXT PRIMARY KEY,
  job      TEXT NOT NULL,
  address  JSONB NOT NULL DEFAULT '[]'::JSONB
);
