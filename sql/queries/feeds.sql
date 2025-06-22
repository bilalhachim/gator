-- name: CreateFeed :one
INSERT INTO feeds(feed_id,created_at,updated_at,name,url,reference_id,last_fetched_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;
-- name: GetFeeds :many
SELECT name,url,reference_id FROM feeds;
-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url=$1;
-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE feed_id = $1;
-- name: GetNextFeedToFetch :one
SELECT  *   FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
