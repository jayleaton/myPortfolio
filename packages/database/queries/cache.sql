-- name: GetCacheEntry :one
SELECT * FROM cache_entry WHERE key = $1;

-- name: PutCacheEntry :exec
INSERT INTO cache_entry (key, value, "updatedAt")
VALUES ($1, $2, $3)
ON CONFLICT (key) DO UPDATE SET value = excluded.value, "updatedAt" = excluded."updatedAt";
