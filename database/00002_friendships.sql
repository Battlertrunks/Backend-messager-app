-- +goose Up
-- +goose StatementBegin
CREATE TABLE friendships (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id          INTEGER references users(id) ON DELETE CASCADE,
  friend_id        INTEGER references users(id) ON DELETE CASCADE,
  status           TEXT NOT NULL,
  created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  deactivated_at   DATETIME
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS friendships;
-- +goose StatementEnd