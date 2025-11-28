-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,
  room_id          INTEGER REFERENCES chat_rooms(id) ON DELETE CASCADE,
  sender_id        INTEGER REFERENCES users(id) ON DELETE SET NULL,
  content          TEXT NOT NULL,
  type             TEXT NOT NULL,
  created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  deactivated_at   DATETIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd