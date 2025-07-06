-- +goose Up
CREATE TABLE playlist (
  `id` CHAR(36) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `slug` VARCHAR(2048) NOT NULL,
  `created_at` TIMESTAMP,
  `updated_at` TIMESTAMP,
  `deleted` CHAR(36),

  PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE playlist;