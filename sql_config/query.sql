-- name: CreateImageSetMeta :execlastid
INSERT INTO img_set_meta (hash, title, origin_url)
VALUES (?, ?, ?);

-- name: CreateImageItem :execlastid
INSERT INTO img_item (ref_meta, hash, data, content_type, url)
VALUES (?, ?, ?, ?, ?);

-- 不支持mysql 批量插入
-- name: CreateBatchImageItem :execresult
-- INSERT INTO img_item (info_ref, data)
-- VALUES ();

-- name: CreateImageItemPV
INSERT INTO img_item_pv (ref_item, data)
VALUES (?, ?);

-- name: GetImageSetMetaID :one
SELECT id
FROM img_set_meta
where title = ?
  and origin_url = ?;

-- name: GetImageSetMetaByHash :one
SELECT id
FROM img_set_meta
WHERE hash = ?;

-- name: GetImageItemByHash :one
SELECT id
FROM img_item
WHERE hash = ?;

-- name: GetImageFailedByHash :one
SELECT id
FROM img_item_failed
WHERE hash = ?;

-- name: CreateImageItemFailed :execlastid
INSERT INTO img_item_failed (ref_meta, hash, failed_url)
VALUES (?, ?, ?);