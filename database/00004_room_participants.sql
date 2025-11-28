-- +goose Up
-- +goose StatementBegin
CREATE TABLE room_participants (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,
  room_id          INTEGER REFERENCES chat_rooms(id) ON DELETE CASCADE,
  user_id          INTEGER REFERENCES users(id) ON DELETE CASCADE,
  created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
  deactivated_at   DATETIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS room_participants;
-- +goose StatementEnd