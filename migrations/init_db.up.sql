CREATE TABLE IF NOT EXISTS url 
(
    id SERIAL PRIMARY KEY,
    alias TEXT UNIQUE,
    url TEXT
);


CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);