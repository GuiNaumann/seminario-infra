CREATE TABLE IF NOT EXISTS books (
    id        SERIAL PRIMARY KEY,
    title     VARCHAR(255) NOT NULL,
    author    VARCHAR(255) NOT NULL,
    year      INT,
    available BOOLEAN DEFAULT TRUE
);

INSERT INTO books (title, author, year, available) VALUES
    ('The Go Programming Language', 'Alan Donovan', 2015, TRUE),
    ('Clean Code', 'Robert C. Martin', 2008, TRUE),
    ('Design Patterns', 'Gang of Four', 1994, FALSE)
ON CONFLICT DO NOTHING;
