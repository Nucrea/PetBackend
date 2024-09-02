create table if not exists shortlinks (
    id text primary key,
    url text,
    expiration date
);