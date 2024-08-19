CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       email VARCHAR(255) NOT NULL UNIQUE,
       name VARCHAR(255) NOT NULL,
       lastname VARCHAR(255) NOT NULL,
       password_hash VARCHAR(255) NOT NULL,
       role VARCHAR(50) NOT NULL,
       image_url TEXT,
       additional_information TEXT,
       course INT,
       major VARCHAR(255),
       degree VARCHAR(255),
       faculty VARCHAR(255),
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);