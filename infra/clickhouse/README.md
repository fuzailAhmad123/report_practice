### activites table
```
CREATE TABLE report_practice.activities
(
    `_id` String,
    `org_id` String,
    `ad_id` String,
    `bets` Float64,
    `wins` Float64,
    `date` DateTime,
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(date)
ORDER BY (org_id, ad_id)
SETTINGS index_granularity = 8192, storage_policy = 'hdd_jbod'
```
