CREATE TABLE IF NOT EXISTS user_ratings(
    id bigint generated always as identity primary key,
    user_id int not null,
    job_id int not null,
    rated_by int not null,
    comment text default null,
    status varchar(1) not null default 'A',
    uuid UUID not null default gen_random_uuid(),
    created_at timestamptz not null default NOW(),
    updated_at timestamptz default NOW(),
    deleted_at timestamptz default null,

    constraint fk_user_ratings_user foreign key (user_id) references users (id),
    constraint fk_user_ratings_job foreign key (job_id) references jobs (id)
)