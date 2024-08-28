CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) not null,
    pwd VARCHAR(255) not null
);

CREATE TABLE IF NOT EXISTS notes
(
    id SERIAL PRIMARY KEY,
    header VARCHAR(255) not null,
    content VARCHAR(255) not null,
    user_id INTEGER REFERENCES users(id)
);