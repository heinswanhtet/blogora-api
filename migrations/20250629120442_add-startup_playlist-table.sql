-- +goose Up
CREATE TABLE startup_playlist (
  `startup_id` CHAR(36) NOT NULL,
  `playlist_id` CHAR(36) NOT NULL,
  `created_at` TIMESTAMP,
  `updated_at` TIMESTAMP,
  `deleted` CHAR(36),

  PRIMARY KEY (startup_id, playlist_id),
  FOREIGN KEY (`startup_id`) REFERENCES startup(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`playlist_id`) REFERENCES playlist(`id`) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE startup_playlist;