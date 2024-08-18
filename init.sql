-- Создание нового пользователя с паролем
CREATE USER myuser WITH PASSWORD 'mypassword';

-- Создание базы данных
CREATE DATABASE mydatabase;

-- Назначение прав доступа к базе данных для пользователя
GRANT ALL PRIVILEGES ON DATABASE mydatabase TO myuser;

-- Даем все права на схемы и таблицы в базе данных
\c mydatabase;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO myuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO myuser;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO myuser;

-- Опционально: делаем пользователя владельцем базы данных
ALTER DATABASE mydatabase OWNER TO myuser;
