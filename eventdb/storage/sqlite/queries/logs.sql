-- name: InsertLogStream :one
insert into log_streams (
	created_at,
	updated_at,
	name,
	token,
	net_whitelist
) values (
	sqlc.arg(created_at),
	sqlc.arg(updated_at),
	sqlc.arg(name),
	sqlc.arg(token),
	sqlc.arg(net_whitelist)
) returning id;

-- name: InsertLogEntry :exec
insert into log_entries (
	stream_id,
	date,
	level,
	message,
	meta
) values (
	sqlc.arg(stream_id),
	sqlc.arg(date),
	sqlc.arg(level),
	sqlc.arg(message),
	sqlc.arg(meta)
);

-- name: QueryLogs :many
select * from log_entries
where (sqlc.narg(stream_id) is null or stream_id = sqlc.narg(stream_id))
order by date
limit ? offset ?;
