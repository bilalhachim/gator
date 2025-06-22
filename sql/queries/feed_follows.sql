-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id,created_at,updated_at,user_id,feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id 
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.feed_id;
-- name: GetFeedFollowsForUser :many
SELECT *
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id 
INNER JOIN feeds ON feed_follows.feed_id = feeds.feed_id
WHERE users.name = $1;

-- name: DeleteFeedFollowsForUser :exec
DELETE FROM feed_follows
WHERE feed_id = (SELECT feed_id FROM feeds WHERE url = $1)
  AND user_id = $2;

