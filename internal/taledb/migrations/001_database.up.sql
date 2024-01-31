CREATE TABLE IF NOT EXISTS tales (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(50),
    url TEXT,
    blurb TEXT,
    added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published TIMESTAMP,
    updated TIMESTAMP,
    review_score SMALLINT,
    review_comment TEXT,
    CONSTRAINT tales_unique_url UNIQUE (url)
    );

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tales_tags (
    tale_id INTEGER REFERENCES tales(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (tale_id, tag_id),
    UNIQUE(tale_id, tag_id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    admin BOOLEAN,
    settings_json JSONB
);

CREATE TABLE IF NOT EXISTS users_tales (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    tale_id INTEGER REFERENCES tales(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tale_id),
    UNIQUE(user_id, tale_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS tales_id ON tales (id);
CREATE UNIQUE INDEX IF NOT EXISTS tales_url ON tales (url);
CREATE UNIQUE INDEX IF NOT EXISTS tags_name ON tags (name);
CREATE UNIQUE INDEX IF NOT EXISTS users_username ON users (username);
CREATE INDEX IF NOT EXISTS tales_tags_tale_id ON tales_tags (tale_id);
CREATE INDEX IF NOT EXISTS users_tales_user_id ON users_tales (user_id);