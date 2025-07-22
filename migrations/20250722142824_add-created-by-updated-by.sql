-- +goose Up

-- author
ALTER TABLE author ADD COLUMN `created_by` CHAR(36);
ALTER TABLE author ADD COLUMN `updated_by` CHAR(36);
ALTER TABLE author ADD CONSTRAINT `author_created_by_fk` FOREIGN KEY(`created_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE author ADD CONSTRAINT `author_updated_by_fk` FOREIGN KEY(`updated_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- startup
ALTER TABLE startup ADD COLUMN `created_by` CHAR(36);
ALTER TABLE startup ADD COLUMN `updated_by` CHAR(36);
ALTER TABLE startup ADD CONSTRAINT `startup_created_by_fk` FOREIGN KEY(`created_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE startup ADD CONSTRAINT `startup_updated_by_fk` FOREIGN KEY(`updated_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- playlist
ALTER TABLE playlist ADD COLUMN `created_by` CHAR(36);
ALTER TABLE playlist ADD COLUMN `updated_by` CHAR(36);
ALTER TABLE playlist ADD CONSTRAINT `playlist_created_by_fk` FOREIGN KEY(`created_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE playlist ADD CONSTRAINT `playlist_updated_by_fk` FOREIGN KEY(`updated_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- startup_playlist
ALTER TABLE startup_playlist ADD COLUMN `created_by` CHAR(36);
ALTER TABLE startup_playlist ADD COLUMN `updated_by` CHAR(36);
ALTER TABLE startup_playlist ADD CONSTRAINT `startup_playlist_created_by_fk` FOREIGN KEY(`created_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE startup_playlist ADD CONSTRAINT `startup_playlist_updated_by_fk` FOREIGN KEY(`updated_by`) REFERENCES author(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;


-- +goose Down

-- author
ALTER TABLE author DROP FOREIGN KEY `author_created_by_fk`;
ALTER TABLE author DROP FOREIGN KEY `author_updated_by_fk`;
ALTER TABLE author DROP COLUMN `created_by`;
ALTER TABLE author DROP COLUMN `updated_by`;

-- startup
ALTER TABLE startup DROP FOREIGN KEY `startup_created_by_fk`;
ALTER TABLE startup DROP FOREIGN KEY `startup_updated_by_fk`;
ALTER TABLE startup DROP COLUMN `created_by`;
ALTER TABLE startup DROP COLUMN `updated_by`;

-- playlist
ALTER TABLE playlist DROP FOREIGN KEY `playlist_created_by_fk`;
ALTER TABLE playlist DROP FOREIGN KEY `playlist_updated_by_fk`;
ALTER TABLE playlist DROP COLUMN `created_by`;
ALTER TABLE playlist DROP COLUMN `updated_by`;

-- startup_playlist
ALTER TABLE startup_playlist DROP FOREIGN KEY `startup_playlist_created_by_fk`;
ALTER TABLE startup_playlist DROP FOREIGN KEY `startup_playlist_updated_by_fk`;
ALTER TABLE startup_playlist DROP COLUMN `created_by`;
ALTER TABLE startup_playlist DROP COLUMN `updated_by`;
