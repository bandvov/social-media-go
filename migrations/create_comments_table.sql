CREATE TYPE entity_type AS ENUM('post', 'comment');
CREATE TYPE status_type AS ENUM('approved', 'rejected');

CREATE TABLE IF NOT EXISTS public.comments ( 
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL,
    entity_id INT NOT NULL,
    entity_type entity_type NOT NULL,
    content TEXT NOT NULL,
    status status_type DEFAULT 'approved',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id)
    );