CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    fisrtname VARCHAR(255) not null,
    lastname VARCHAR(255) not null,
    username VARCHAR(255) not null,
    pwd VARCHAR(255) not null
);

CREATE TABLE IF NOT EXISTS notes
(
    id SERIAL PRIMARY KEY,
    header VARCHAR(255) not null,
    content VARCHAR(10000) not null,
    user_id INTEGER REFERENCES users(id)
);

INSERT INTO users(fisrtname, lastname, username, pwd) values 
('Terry', 'Davis', 'tdavis', 'c3VwZXJwYXNzd29yZA=='),
('Ryan', 'Gosling', 'rgosling', 'ZHJpdmU='),
('Richard', 'Stallman', 'rstallman', 'Z251bGludXg=');
-- superpassword
-- drive
-- gnulinux
