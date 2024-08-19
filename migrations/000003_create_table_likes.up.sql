CREATE TABLE likes (
       id SERIAL PRIMARY KEY,
       user_id INT NOT NULL,
       entity_id INT NOT NULL,
       entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('post', 'comment', 'topic')),
       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
       UNIQUE (user_id, entity_id, entity_type)
);


ALTER TABLE likes ADD CONSTRAINT fk_like_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Индекс для ускорения запросов по entity_id и entity_type
CREATE INDEX idx_likes_entity ON likes(entity_id, entity_type);
