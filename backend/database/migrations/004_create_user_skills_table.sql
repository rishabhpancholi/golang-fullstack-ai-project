-- Write your migrate up statements here

CREATE TABLE user_skills (
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    skill_id INTEGER NOT NULL REFERENCES skills(skill_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, skill_id)
);

---- create above / drop below ----

DROP TABLE user_skills;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
