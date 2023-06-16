CREATE TABLE IF NOT EXISTS "user" (
    id serial PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    role VARCHAR(255)
);