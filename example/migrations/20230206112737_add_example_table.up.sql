CREATE TABLE IF NOT EXISTS example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO example ("id", "name") VALUES (1, 'foo 1');
INSERT INTO example ("id", "name") VALUES (2, 'foo 2');

CREATE TABLE IF NOT EXISTS example_2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    example_1_id INT NOT NULL REFERENCES example(id)
);
