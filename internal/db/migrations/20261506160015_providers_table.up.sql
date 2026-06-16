CREATE TABLE IF NOT EXISTS providers(
    id bigint generated always as identity primary key,
    name varchar(255) not null unique,
    provider_id varchar(45) not null unique,
    addr1 varchar(255),
    addr2 varchar(255),
    city varchar(255),
    postal varchar(255),
    latitude numeric(9,6),
    longitude numeric(10,6),
    status varchar(1) not null default 'A',
    uuid UUID not null default gen_random_uuid(),
    created_at timestamptz not null default NOW(),
    updated_at timestamptz default NOW(),
    deleted_at timestamptz default null
)