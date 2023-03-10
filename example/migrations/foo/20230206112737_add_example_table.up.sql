CREATE TABLE IF NOT EXISTS foo.example (
    id SERIAL PRIMARY KEY,
    example_schema_id INT NOT NULL REFERENCES bar.example(id),
    name VARCHAR(255) NOT NULL
);

INSERT INTO foo.example ("id", "name", "example_schema_id") VALUES (1, 'foo 1', 88);
INSERT INTO foo.example ("id", "name", "example_schema_id") VALUES (2, 'foo 2', 99);

CREATE TABLE IF NOT EXISTS foo.example_2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    example_1_id INT NOT NULL REFERENCES foo.example(id)
);
