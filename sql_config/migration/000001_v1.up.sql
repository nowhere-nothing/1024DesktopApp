CREATE TABLE img_set_meta
(
    id         BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title      TEXT   NOT NULL,
    origin_url TEXT   NOT NULL
);

CREATE TABLE img_item
(
    id       BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ref_meta BIGINT NOT NULL,
    data     MEDIUMBLOB
);

CREATE TABLE img_item_pv
(
    ref_item BIGINT NOT NULL PRIMARY KEY,
    data     MEDIUMBLOB
);

CREATE TABLE img_item_failed
(
    id         BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ref_meta   BIGINT NOT NULL,
    failed_url TEXT   NOT NULL
)
