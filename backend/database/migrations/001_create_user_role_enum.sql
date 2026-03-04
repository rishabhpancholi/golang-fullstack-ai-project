-- Write your migrate up statements here

CREATE TYPE user_role AS ENUM ('jobseeker', 'recruiter', 'admin');

---- create above / drop below ----

DROP TYPE user_role;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
