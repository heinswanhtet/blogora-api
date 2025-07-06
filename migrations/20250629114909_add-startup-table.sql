-- +goose Up
CREATE TABLE startup (
  `id` CHAR(36) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `slug` VARCHAR(2048) NOT NULL,
  `author_id` CHAR(36) NOT NULL,
  `views` INT DEFAULT 0,
  `description` TEXT,
  `category` VARCHAR(255) NOT NULL,
  `image` VARCHAR(2048),
  `pitch` TEXT,
  `created_at` TIMESTAMP,
  `updated_at` TIMESTAMP,
  `deleted` CHAR(36),

  PRIMARY KEY (id),
  FOREIGN KEY (`author_id`) REFERENCES author(`id`) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE startup;