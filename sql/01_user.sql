create table if not exists users (
    id int generated always as identity,
    email text unique not null,
    secret text not null,
    name text not null,
    email_verified boolean not null default false,

    primary key (id)
);

create index if not exists users_email_idx on users(email);