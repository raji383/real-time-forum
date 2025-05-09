CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT NOT NULL,
    age INTEGER,
    gender TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE,
    password TEXT
);
No posts found