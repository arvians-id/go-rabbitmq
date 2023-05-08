CREATE TABLE IF NOT EXISTS category_todo(
    todo_id INT REFERENCES todos(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY(todo_id, category_id)
);
