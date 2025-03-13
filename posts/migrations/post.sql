CREATE TABLE
    posts (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW ()
    );

CREATE TABLE
    tags (id SERIAL PRIMARY KEY, tag TEXT NOT NULL UNIQUE);

CREATE TABLE
    keywords (
        id SERIAL PRIMARY KEY,
        keyword TEXT NOT NULL UNIQUE
    );

CREATE TABLE
    post_tags (
        post_id INT REFERENCES posts (id) ON DELETE CASCADE,
        tag_id INT REFERENCES tags (id) ON DELETE CASCADE,
        PRIMARY KEY (post_id, tag_id)
    );

CREATE TABLE
    post_keywords (
        post_id INT REFERENCES posts (id) ON DELETE CASCADE,
        keyword_id INT REFERENCES keywords (id) ON DELETE CASCADE,
        PRIMARY KEY (post_id, keyword_id)
    );