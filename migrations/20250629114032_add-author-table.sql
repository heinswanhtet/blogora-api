-- +goose Up
CREATE TABLE author (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `username` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `image` VARCHAR(2048),
  `bio` TEXT,
  `created_at` TIMESTAMP,
  `updated_at` TIMESTAMP,
  `deleted` CHAR(36),

  PRIMARY KEY (id),
  UNIQUE KEY(email)
);

-- +goose Down
DROP TABLE author;