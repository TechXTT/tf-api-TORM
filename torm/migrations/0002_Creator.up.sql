CREATE TABLE creator (
    id SERIAL PRIMARY KEY,
    createdat TIMESTAMP DEFAULT now(),
    updatedat TIMESTAMP DEFAULT now(),
    deletedat TIMESTAMP,
    name TEXT,
    email TEXT,
    phone TEXT,
    grade INTEGER,
    class TEXT
);

CREATE INDEX idx_creator_1 ON creator (grade);

CREATE INDEX idx_creator_2 ON creator (class);