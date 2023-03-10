CREATE TABLE IF NOT EXISTS bar.example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO bar.example ("id", "name") VALUES (88, 'foo 1');
INSERT INTO bar.example ("id", "name") VALUES (99, 'foo 2');

CREATE TABLE IF NOT EXISTS bar.example_2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    example_1_id INT NOT NULL REFERENCES bar.example(id)
);
