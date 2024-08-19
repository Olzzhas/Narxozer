CREATE TABLE IF NOT EXISTS clubs (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  image_url TEXT,
    creator INT not null references users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  date TIMESTAMPTZ NOT NULL,
  club_id INTEGER NOT NULL REFERENCES clubs(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS club_members (
  club_id INTEGER NOT NULL REFERENCES clubs(id),
  user_id INTEGER NOT NULL REFERENCES users(id),
  PRIMARY KEY (club_id, user_id)
);

CREATE TABLE IF NOT EXISTS club_admins (
  club_id INTEGER NOT NULL REFERENCES clubs(id),
  user_id INTEGER NOT NULL REFERENCES users(id),
  PRIMARY KEY (club_id, user_id)
);
