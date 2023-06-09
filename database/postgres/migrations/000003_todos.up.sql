CREATE TABLE IF NOT EXISTS todos(
    id serial PRIMARY KEY,
    title VARCHAR (50) NOT NULL,
    description TEXT NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
