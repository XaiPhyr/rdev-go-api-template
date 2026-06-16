CREATE TABLE IF NOT EXISTS service_types(
    id bigint generated always as identity primary key,
    name varchar(255) not null unique,
    description text default null,
    price bigint default null,
    status varchar(1) not null default 'A',
    uuid UUID not null default gen_random_uuid(),
    created_at timestamptz not null default NOW(),
    updated_at timestamptz default NOW(),
    deleted_at timestamptz default null
)