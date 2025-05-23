begin;

create table log_streams (
	id blob primary key,
	created_at integer not null,
	updated_at integer not null,
	name text not null unique,
	token blob,
	platform text null,
	net_whitelist blob
);

create table log_entries (
	id integer primary key autoincrement,
	stream_id blob not null references log_streams(id) on update restrict on delete cascade,
	date integer not null,
	level text not null,
	message text not null,
	meta blob
);

end;
