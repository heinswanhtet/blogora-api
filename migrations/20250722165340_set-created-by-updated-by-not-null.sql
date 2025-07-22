-- +goose Up

-- after setting created_by and updated_by in old data

-- author
ALTER TABLE author CHANGE `created_by` `created_by` char(36) NOT NULL;
ALTER TABLE author CHANGE `updated_by` `updated_by` char(36) NOT NULL;

-- startup
ALTER TABLE startup CHANGE `created_by` `created_by` char(36) NOT NULL;
ALTER TABLE startup CHANGE `updated_by` `updated_by` char(36) NOT NULL;

-- playlist
ALTER TABLE playlist CHANGE `created_by` `created_by` char(36) NOT NULL;
ALTER TABLE playlist CHANGE `updated_by` `updated_by` char(36) NOT NULL;

-- startup_playlist
ALTER TABLE startup_playlist CHANGE `created_by` `created_by` char(36) NOT NULL;
ALTER TABLE startup_playlist CHANGE `updated_by` `updated_by` char(36) NOT NULL;


-- +goose Down

-- author
ALTER TABLE author CHANGE `created_by` `created_by` char(36) NULL;
ALTER TABLE author CHANGE `updated_by` `updated_by` char(36) NULL;

-- startup
ALTER TABLE startup CHANGE `created_by` `created_by` char(36) NULL;
ALTER TABLE startup CHANGE `updated_by` `updated_by` char(36) NULL;

-- playlist
ALTER TABLE playlist CHANGE `created_by` `created_by` char(36) NULL;
ALTER TABLE playlist CHANGE `updated_by` `updated_by` char(36) NULL;

-- startup_playlist
ALTER TABLE startup_playlist CHANGE `created_by` `created_by` char(36) NULL;
ALTER TABLE startup_playlist CHANGE `updated_by` `updated_by` char(36) NULL;