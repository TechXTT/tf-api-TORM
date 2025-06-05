CREATE TABLE vote (
    id SERIAL PRIMARY KEY,
    createdat TIMESTAMP DEFAULT now(),
    updatedat TIMESTAMP DEFAULT now(),
    deletedat TIMESTAMP,
    name TEXT,
    email TEXT,
    verified BOOLEAN DEFAULT false,
    networksid INTEGER,
    softwareid INTEGER,
    embeddedid INTEGER,
    battlebotid INTEGER
);

CREATE INDEX idx_vote_1 ON vote (networksid);

CREATE INDEX idx_vote_2 ON vote (softwareid);

CREATE INDEX idx_vote_3 ON vote (embeddedid);

CREATE INDEX idx_vote_4 ON vote (battlebotid);