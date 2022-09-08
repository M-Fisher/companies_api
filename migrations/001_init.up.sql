CREATE TABLE IF NOT EXISTS companies(
    id       serial PRIMARY KEY,
    name     TEXT,
    code     TEXT,
    country  TEXT,
    website  TEXT      DEFAULT NULL,
    phone    TEXT
);