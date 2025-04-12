-- name: GetDelegationSnapshotByValidatorAndDelegator :one
 SELECT id, validator_address, 
        delegator_address, amount_uatom, 
        change_uatom, timestamp 
    FROM delegation_snapshots
    WHERE validator_address = $1 AND delegator_address = $2
    ORDER BY timestamp DESC LIMIT 1;

-- name: GetDelegationSnapshotByValidator :many
 SELECT delegator_address, amount_uatom, timestamp, change_uatom
    FROM delegation_snapshots
    WHERE validator_address = $1
    ORDER BY timestamp ASC
    LIMIT $2
    OFFSET $3;

-- name: GetCountDelegationSnapshotByValidator :one
 SELECT COUNT(*)
    FROM delegation_snapshots
    WHERE validator_address = $1;

-- name: GetDailyAggregateByValidator :many
 SELECT delegator_address, date, total_amount
    FROM daily_aggregates
    WHERE validator_address = $1
    ORDER BY date ASC
    LIMIT $2
    OFFSET $3;

-- name: GetCountDailyAggregateByValidator :one
 SELECT COUNT(*)
    FROM daily_aggregates
    WHERE validator_address = $1;

-- name: GetDelegatorHistoryByValidator :many
SELECT timestamp, amount_uatom, change_uatom
    FROM delegation_snapshots
    WHERE validator_address = $1 AND delegator_address = $2
    ORDER BY 
    CASE WHEN @sort_by::text = '-date' THEN "timestamp" END DESC,
    CASE WHEN @sort_by::text = 'date' THEN "timestamp" END ASC
    LIMIT $3
    OFFSET $4;

-- name: GetCountDelegatorHistoryByValidator :one
SELECT COUNT(*)
    FROM delegation_snapshots
    WHERE validator_address = $1 AND delegator_address = $2;

-- name: CreateDelegationSnapshot :one
INSERT INTO delegation_snapshots (
    validator_address,
    delegator_address,
    amount_uatom,
    change_uatom,
    timestamp
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetLatestDelegationSnapshot :many
SELECT DISTINCT ON (delegator_address, validator_address)
           validator_address, delegator_address, amount_uatom
    FROM delegation_snapshots
    ORDER BY delegator_address, validator_address, timestamp DESC;

-- name: CreateDailyAggregate :one
INSERT INTO daily_aggregates (validator_address, delegator_address, date, total_amount)
VALUES ($1, $2, $3, $4) RETURNING id;