CREATE TABLE project (
    id SERIAL PRIMARY KEY,
    createdat TIMESTAMP DEFAULT now(),
    updatedat TIMESTAMP DEFAULT now(),
    deletedat TIMESTAMP,
    name TEXT,
    description TEXT,
    type TEXT,
    category TEXT,
    mentor TEXT,
    videolink TEXT,
    hasthumbnail BOOLEAN DEFAULT false,
    demolink TEXT,
    githublink TEXT
);

CREATE INDEX idx_project_1 ON project (type);

CREATE INDEX idx_project_2 ON project (category);