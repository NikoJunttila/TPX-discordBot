-- +goose Up 

CREATE TABLE usersCount (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  count INT NOT NULL
);

-- +goose Down
DROP TABLE usersCount;