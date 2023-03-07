ALTER TABLE img_set_meta
    ADD COLUMN hash VARCHAR(256) NOT NULL UNIQUE;
ALTER TABLE img_set_meta
    ADD INDEX idx_hash (hash);

ALTER TABLE img_item
    ADD COLUMN hash VARCHAR(256) NOT NULL UNIQUE;
ALTER TABLE img_item
    ADD INDEX idx_hash (hash);

ALTER TABLE img_item_failed
    ADD COLUMN hash VARCHAR(256) NOT NULL UNIQUE;
ALTER TABLE img_item_failed
    ADD INDEX idx_hash (hash);
