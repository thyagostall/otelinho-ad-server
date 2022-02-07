create table if not exists campaigns
(
	id serial not null
		constraint campaigns_pk
			primary key,
	creative varchar(100) not null,
	start_date timestamp not null,
	end_date timestamp not null,
	goal integer not null,
	max_bid numeric not null,
	budget numeric
);

alter table campaigns owner to otelinho;

create table if not exists beacons
(
	id serial not null
		constraint beacons_pk
			primary key,
	campaign_id integer not null
		constraint beacons_campaigns_fk
			references campaigns,
	event varchar(100) not null,
	timestamp timestamp default now() not null
);

alter table beacons owner to otelinho;

create table if not exists pacing
(
	campaign_id integer not null
		constraint pacing_pk
			primary key
		constraint pacing_campaign_fk
			references campaigns,
	velocity bigint not null
);

alter table pacing owner to otelinho;

create unique index if not exists pacing_campaign_id_uindex
	on pacing (campaign_id);

create table if not exists ad_requests
(
	id serial not null
		constraint ad_requests_pk
			primary key,
	timestamp timestamp default now() not null
);

alter table ad_requests owner to otelinho;
