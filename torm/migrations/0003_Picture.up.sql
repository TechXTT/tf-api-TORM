CREATE TABLE picture (
    id SERIAL PRIMARY KEY,
    createdat TIMESTAMP DEFAULT now(),
    updatedat TIMESTAMP DEFAULT now(),
    deletedat TIMESTAMP,
    url TEXT,
    isthumbnail BOOLEAN DEFAULT false,
    projectid INTEGER
);

CREATE INDEX idx_picture_1 ON picture (projectid);