-- +goose Up
-- Last good values of the site's live data (X followers, latest YouTube
-- videos), so a restart serves them before the first refresh lands.
CREATE TABLE cache_entry (
    key          text PRIMARY KEY,
    value        jsonb NOT NULL,
    "updatedAt"  timestamptz NOT NULL
);

-- +goose Down
DROP TABLE cache_entry;
