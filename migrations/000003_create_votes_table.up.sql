CREATE TABLE votes (
    user_id INTEGER REFERENCES users(id),
    post_id INTEGER REFERENCES posts(id),
    vote INTEGER,
    PRIMARY KEY (user_id, post_id),
    CONSTRAINT check_vote_value CHECK (vote IN (-1, 1))
);
