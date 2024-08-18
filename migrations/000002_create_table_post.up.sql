CREATE TABLE posts (
       id SERIAL PRIMARY KEY,
       title TEXT NOT NULL,
       content TEXT NOT NULL,
       image_url TEXT,
       author_id INT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
       updated_at TIMESTAMPTZ,
       likes INT DEFAULT 0
);

CREATE TABLE comments (
      id SERIAL PRIMARY KEY,
      post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
      content TEXT NOT NULL,
      author_id INT NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
      likes INT DEFAULT 0,
      parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE
);


ALTER TABLE posts ADD CONSTRAINT fk_post_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE comments ADD CONSTRAINT fk_comment_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;
