CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       email TEXT NOT NULL UNIQUE,
       name TEXT NOT NULL,
       lastname TEXT NOT NULL,
       password_hash TEXT NOT NULL,
       role VARCHAR(10) NOT NULL CHECK (role IN ('STUDENT', 'TEACHER', 'ADMIN')),
       image_url TEXT,
       additional_information TEXT,
       course INT CHECK (course >= 1 AND course <= 4),
       major TEXT,
       degree TEXT,
       faculty TEXT,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
       updated_at TIMESTAMPTZ
);
