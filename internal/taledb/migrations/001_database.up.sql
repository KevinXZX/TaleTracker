CREATE TABLE IF NOT EXISTS tales (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    url TEXT,
    blurb TEXT,
    added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published TIMESTAMP,
    updated TIMESTAMP,
    review_score SMALLINT,
    review_comment TEXT,
    CONSTRAINT tales_unique_url UNIQUE (url)
);


CREATE UNIQUE INDEX IF NOT EXISTS tales_id ON tales (id);
CREATE UNIQUE INDEX IF NOT EXISTS tales_url ON tales (url);