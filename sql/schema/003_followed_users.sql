-- +goose Up 
CREATE TABLE follow_users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  account_name TEXT NOT NULL,
  hashtag TEXT NOT NULL,
  puuID TEXT NOT NULL UNIQUE,
  region TEXT NOT NULL
);

-- +goose Down
DROP TABLE follow_users;