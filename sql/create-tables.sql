CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    release_date DATE NOT NULL,
    rating FLOAT CHECK (rating >= 0 AND rating <= 10),
    actors_id INT[] NOT NULL,
    deleted_at DATE
);

CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sex VARCHAR(10) NOT NULL,
    birthday DATE NOT NULL,
    movies_id INT[],
    deleted_at DATE
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(30) UNIQUE NOT NULL,
    role VARCHAR(10) NOT NULL,
    password_hash VARCHAR(60) NOT NULL
);