create table if not exists shortlinks (
    id text primary key,
    url text not null,
    expires_at timestamp not null
);
