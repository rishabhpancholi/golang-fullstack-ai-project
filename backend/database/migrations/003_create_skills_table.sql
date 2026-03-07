-- Write your migrate up statements here

CREATE TABLE skills (
    skill_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

---- create above / drop below ----

DROP TABLE skills;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
